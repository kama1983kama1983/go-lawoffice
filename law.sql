CREATE TABLE Clients (
	    client_id INT AUTO_INCREMENT PRIMARY KEY,
	    first_name VARCHAR(100),
	    last_name VARCHAR(100),
	    email VARCHAR(255),
	    phone VARCHAR(20),
	    address VARCHAR(255),
	    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);
	
	CREATE TABLE `Lawyers` (
	`lawyer_id` int(11) NOT NULL AUTO_INCREMENT,
	`first_name` varchar(100) DEFAULT NULL,
	`last_name` varchar(100) DEFAULT NULL,
	`email` varchar(255) DEFAULT NULL,
	`phone` varchar(20) DEFAULT NULL,
	`specialization` varchar(200) DEFAULT NULL,
	`number_lawyer` varchar(200) DEFAULT NULL,
	`profice_pic` varchar(200) DEFAULT NULL,
	`updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	`created_at` datetime DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`lawyer_id`)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;

	CREATE TABLE PowerOfAttorney (
    power_of_attorney_id INT AUTO_INCREMENT PRIMARY KEY,
    client_id INT,
    lawyer_id INT,
    case_id INT,
    date_of_issue DATETIME DEFAULT CURRENT_TIMESTAMP,
    expiration_date DATETIME,
    description TEXT,
    reference_code VARCHAR(10), -- حقل جديد لتخزين الرقم والحرف والسنة
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (client_id) REFERENCES Clients(client_id),
    FOREIGN KEY (lawyer_id) REFERENCES Lawyers(lawyer_id),
    FOREIGN KEY (case_id) REFERENCES Cases(case_id)
);

DROP TABLE IF EXISTS `Sessions`;
CREATE TABLE `Sessions` (
  `Session_id` int(11) NOT NULL AUTO_INCREMENT,
  `Case_id` int(11) DEFAULT NULL,
  `qrar` varchar(100) DEFAULT NULL,
  `today` varchar(100) DEFAULT NULL,
  `next_date` varchar(255) DEFAULT NULL,
  `description` text,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`Session_id`),
  KEY `Case_id` (`Case_id`),
  CONSTRAINT `Sessions_ibfk_1` FOREIGN KEY (`Case_id`) REFERENCES `Cases` (`case_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	
	CREATE TABLE Lawyers (
	    lawyer_id INT AUTO_INCREMENT PRIMARY KEY,
	    first_name VARCHAR(100),
	    last_name VARCHAR(100),
	    email VARCHAR(255),
	    phone VARCHAR(20),
	    specialization VARCHAR(100),
	    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);
	DROP TABLE IF EXISTS `Cases`;
CREATE TABLE `Cases` (
  `case_id` int(11) NOT NULL AUTO_INCREMENT,
  `client_id` int(11) DEFAULT NULL,
  `lawyer_id` int(11) DEFAULT NULL,
  `case_number` varchar(50) DEFAULT NULL,
  `case_type` varchar(100) DEFAULT NULL,
  `status` enum('Open','Closed','Pending','On Hold') DEFAULT NULL,
  `description` text,
  `court` varchar(100) DEFAULT NULL,       -- حقل جديد لمحكمة
  `dayra` varchar(100) DEFAULT NULL,       -- حقل جديد لدائرة
  `subject` varchar(255) DEFAULT NULL,     -- حقل جديد لموضوع القضية
  `mod3le` varchar(100) DEFAULT NULL,      -- حقل جديد لنموذج القضية
  `ed3a_date` date DEFAULT NULL,           -- حقل جديد لتاريخ الإدعا
  `qa3a_num` varchar(50) DEFAULT NULL,     -- حقل جديد لرقم القاعة
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`case_id`),
  KEY `client_id` (`client_id`),
  KEY `lawyer_id` (`lawyer_id`),
  CONSTRAINT `Cases_ibfk_1` FOREIGN KEY (`client_id`) REFERENCES `Clients` (`client_id`),
  CONSTRAINT `Cases_ibfk_2` FOREIGN KEY (`lawyer_id`) REFERENCES `Lawyers` (`lawyer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	
	CREATE TABLE Appointments (
	    appointment_id INT AUTO_INCREMENT PRIMARY KEY,
	    client_id INT,
	    lawyer_id INT,
	    appointment_date DATETIME,
	    location VARCHAR(255),
	    notes TEXT,
	    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	    FOREIGN KEY (client_id) REFERENCES Clients(client_id),
	    FOREIGN KEY (lawyer_id) REFERENCES Lawyers(lawyer_id)
	);
	
	CREATE TABLE Invoices (
	    invoice_id INT AUTO_INCREMENT PRIMARY KEY,
	    client_id INT,
	    case_id INT,
	    amount DECIMAL(10, 2),
	    status ENUM('Paid', 'Unpaid', 'Overdue'),
	   	due_date VARCHAR(255) ,
	    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	    FOREIGN KEY (case_id) REFERENCES Cases(case_id),
	    FOREIGN KEY (client_id) REFERENCES Clients(client_id)
	);


-- Creating the Tasks table

CREATE TABLE Tasks (

    task_id INT AUTO_INCREMENT PRIMARY KEY,

    session_id INT,

    task_description TEXT,

    reminder_date DATETIME,

    notified TINYINT DEFAULT 0, -- 0 for not notified, 1 for notified

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (session_id) REFERENCES Sessions(Session_id)

);


-- Example Trigger to create a task reminder for sessions 15 days before the next_date

DELIMITER //

CREATE TRIGGER create_reminder_before_session

AFTER INSERT ON Sessions

FOR EACH ROW

BEGIN

    DECLARE reminder_date DATETIME;

    SET reminder_date = DATE_SUB(NEW.next_date, INTERVAL 15 DAY);

    

    INSERT INTO Tasks (session_id, task_description, reminder_date)

    VALUES (NEW.Session_id, CONCAT('Reminder for session on ', NEW.next_date), reminder_date);

END; //

DELIMITER ;


-- Example Trigger to set a notified status for tasks 7 days after the reminder date

DELIMITER //

CREATE TRIGGER notify_task_after_reminder

AFTER INSERT ON Tasks

FOR EACH ROW

BEGIN

    DECLARE notify_date DATETIME;

    SET notify_date = DATE_ADD(NEW.reminder_date, INTERVAL 7 DAY);

    

    -- Assuming a separate process or application logic will handle the notification

    -- You can log or handle notification logic here

END; //

DELIMITER ;