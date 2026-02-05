-- PT XYZ Multifinance Database Schema
-- ACID Compliant Design with proper indexes and constraints

-- Create database if not exists
CREATE DATABASE IF NOT EXISTS xyz_multifinance;
USE xyz_multifinance;

-- Drop existing tables (if any)
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS consumer_limits;
DROP TABLE IF EXISTS consumers;

-- Table: Consumers (Customer Data)
-- Stores personal information of PT XYZ Multifinance customers
CREATE TABLE consumers (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    nik VARCHAR(16) NOT NULL UNIQUE COMMENT 'Nomor KTP Konsumen',
    full_name VARCHAR(255) NOT NULL COMMENT 'Nama Lengkap Konsumen',
    legal_name VARCHAR(255) NOT NULL COMMENT 'Nama Resmi di KTP',
    place_of_birth VARCHAR(255) COMMENT 'Tempat Lahir',
    date_of_birth DATETIME COMMENT 'Tanggal Lahir',
    salary DECIMAL(15, 2) NOT NULL COMMENT 'Gaji Konsumen',
    ktp_photo LONGTEXT COMMENT 'Foto KTP (Base64)',
    selfie_photo LONGTEXT COMMENT 'Foto Selfie Konsumen (Base64)',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    INDEX idx_nik (nik),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT check_salary CHECK (salary >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Tabel Konsumen PT XYZ Multifinance';

-- Table: Consumer Limits (Credit Limits)
-- Tracks credit limits per tenor (1, 2, 3, 6 months) for each consumer
CREATE TABLE consumer_limits (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    consumer_id BIGINT UNSIGNED NOT NULL,
    tenor INT NOT NULL COMMENT 'Tenor in months: 1, 2, 3, 6',
    limit_amount DECIMAL(15, 2) NOT NULL COMMENT 'Credit limit amount',
    used_amount DECIMAL(15, 2) DEFAULT 0 COMMENT 'Used amount from limit',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (consumer_id) REFERENCES consumers(id) ON DELETE CASCADE ON UPDATE CASCADE,
    UNIQUE KEY unique_consumer_tenor (consumer_id, tenor),
    INDEX idx_consumer_id (consumer_id),
    INDEX idx_tenor (tenor),
    CONSTRAINT check_tenor CHECK (tenor IN (1, 2, 3, 6)),
    CONSTRAINT check_limit_amount CHECK (limit_amount > 0),
    CONSTRAINT check_used_amount CHECK (used_amount >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Tabel Limit Kredit Konsumen';

-- Table: Transactions (Financial Transactions)
-- Records all financing transactions (purchases with installments)
CREATE TABLE transactions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    consumer_id BIGINT UNSIGNED NOT NULL,
    contract_number VARCHAR(255) NOT NULL UNIQUE COMMENT 'Nomor Kontrak Unik',
    tenor INT NOT NULL COMMENT 'Tenor in months: 1, 2, 3, 6',
    otr DECIMAL(15, 2) NOT NULL COMMENT 'On The Road (OTR) Price',
    admin_fee DECIMAL(15, 2) DEFAULT 0 COMMENT 'Admin Fee',
    installment_amount DECIMAL(15, 2) NOT NULL COMMENT 'Jumlah Cicilan',
    interest_amount DECIMAL(15, 2) DEFAULT 0 COMMENT 'Bunga',
    asset_name VARCHAR(255) COMMENT 'Nama Aset yang Dibeli',
    status VARCHAR(50) DEFAULT 'ACTIVE' COMMENT 'ACTIVE, COMPLETED, DEFAULTED',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (consumer_id) REFERENCES consumers(id) ON DELETE CASCADE ON UPDATE CASCADE,
    UNIQUE KEY unique_contract_number (contract_number),
    INDEX idx_consumer_id (consumer_id),
    INDEX idx_contract_number (contract_number),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_tenor (tenor),
    CONSTRAINT check_otr CHECK (otr > 0),
    CONSTRAINT check_installment CHECK (installment_amount > 0),
    CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'DEFAULTED'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Tabel Transaksi Pembiayaan';

-- Create views for business intelligence

-- View: Consumer Overview
CREATE VIEW v_consumer_overview AS
SELECT 
    c.id,
    c.nik,
    c.full_name,
    c.salary,
    COUNT(DISTINCT cl.id) as total_limits,
    SUM(cl.limit_amount) as total_limit_amount,
    SUM(cl.used_amount) as total_used_amount,
    COUNT(DISTINCT t.id) as total_transactions
FROM consumers c
LEFT JOIN consumer_limits cl ON c.id = cl.consumer_id
LEFT JOIN transactions t ON c.id = t.consumer_id AND t.status = 'ACTIVE'
WHERE c.deleted_at IS NULL
GROUP BY c.id, c.nik, c.full_name, c.salary;

-- View: Transaction Summary
CREATE VIEW v_transaction_summary AS
SELECT 
    t.id,
    t.consumer_id,
    c.full_name,
    t.contract_number,
    t.tenor,
    t.otr,
    t.admin_fee,
    t.installment_amount,
    t.interest_amount,
    t.asset_name,
    t.status,
    (t.otr + t.admin_fee + (t.interest_amount * t.tenor)) as total_amount,
    t.created_at
FROM transactions t
INNER JOIN consumers c ON t.consumer_id = c.id
WHERE c.deleted_at IS NULL;

-- Indexes for performance optimization

-- Additional index for transaction queries
CREATE INDEX idx_transaction_consumer_status ON transactions(consumer_id, status);

-- Index for limit utilization queries
CREATE INDEX idx_limit_utilization ON consumer_limits(consumer_id, used_amount, limit_amount);

-- Index for date range queries
CREATE INDEX idx_transaction_date_range ON transactions(created_at, consumer_id);

-- Stored Procedures for ACID Compliance

-- Procedure: Process New Transaction (Atomic operation)
DELIMITER $$

CREATE PROCEDURE sp_create_transaction(
    IN p_consumer_id BIGINT UNSIGNED,
    IN p_contract_number VARCHAR(255),
    IN p_tenor INT,
    IN p_otr DECIMAL(15, 2),
    IN p_admin_fee DECIMAL(15, 2),
    IN p_installment_amount DECIMAL(15, 2),
    IN p_interest_amount DECIMAL(15, 2),
    IN p_asset_name VARCHAR(255),
    OUT p_transaction_id BIGINT UNSIGNED,
    OUT p_error_message VARCHAR(500)
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        SET p_error_message = 'Transaction failed - Database error';
        ROLLBACK;
    END;

    START TRANSACTION;

    -- Check consumer exists
    IF NOT EXISTS (SELECT 1 FROM consumers WHERE id = p_consumer_id AND deleted_at IS NULL) THEN
        SET p_error_message = 'Consumer not found';
        ROLLBACK;
        LEAVE sp_create_transaction;
    END IF;

    -- Check contract number uniqueness
    IF EXISTS (SELECT 1 FROM transactions WHERE contract_number = p_contract_number) THEN
        SET p_error_message = 'Contract number already exists';
        ROLLBACK;
        LEAVE sp_create_transaction;
    END IF;

    -- Check and update limit
    IF NOT EXISTS (
        SELECT 1 FROM consumer_limits 
        WHERE consumer_id = p_consumer_id 
        AND tenor = p_tenor 
        AND (used_amount + p_otr) <= limit_amount
    ) THEN
        SET p_error_message = 'Insufficient credit limit';
        ROLLBACK;
        LEAVE sp_create_transaction;
    END IF;

    -- Update limit usage
    UPDATE consumer_limits 
    SET used_amount = used_amount + p_otr,
        updated_at = CURRENT_TIMESTAMP
    WHERE consumer_id = p_consumer_id AND tenor = p_tenor;

    -- Insert transaction
    INSERT INTO transactions (
        consumer_id, contract_number, tenor, otr, admin_fee, 
        installment_amount, interest_amount, asset_name, status
    ) VALUES (
        p_consumer_id, p_contract_number, p_tenor, p_otr, p_admin_fee,
        p_installment_amount, p_interest_amount, p_asset_name, 'ACTIVE'
    );

    SET p_transaction_id = LAST_INSERT_ID();
    SET p_error_message = NULL;

    COMMIT;
END$$

DELIMITER ;

-- Stored Procedure: Get Consumer Available Limit
DELIMITER $$

CREATE PROCEDURE sp_get_consumer_available_limit(
    IN p_consumer_id BIGINT UNSIGNED,
    IN p_tenor INT,
    OUT p_available_limit DECIMAL(15, 2)
)
BEGIN
    SELECT (limit_amount - used_amount) INTO p_available_limit
    FROM consumer_limits
    WHERE consumer_id = p_consumer_id AND tenor = p_tenor;
END$$

DELIMITER ;

-- Set session variables for transaction isolation
SET SESSION TRANSACTION ISOLATION LEVEL READ_COMMITTED;

COMMIT;
