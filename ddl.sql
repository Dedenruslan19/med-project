CREATE DATABASE p2;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    weight DECIMAL(5,2) NOT NULL,
    height DECIMAL(5,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE workouts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    goals VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE exercises (
    id SERIAL PRIMARY KEY,
    workout_id INTEGER NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    sets VARCHAR(255),
    reps VARCHAR(255),
    equipment VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE exercise_logs (
    id SERIAL PRIMARY KEY,
    exercise_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    set_count INTEGER NOT NULL,
    rep_count INTEGER NOT NULL,
    weight DECIMAL(5,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE doctors (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    specialization VARCHAR(255) NOT NULL,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE appointments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    doctor_id INTEGER NOT NULL,
    appointment_date TIMESTAMP NOT NULL,
    status VARCHAR(50) DEFAULT 'pending', 
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE
);

CREATE TABLE diagnoses (
    id SERIAL PRIMARY KEY,
    appointment_id INTEGER NOT NULL UNIQUE,
    doctor_id INTEGER NOT NULL,
    notes TEXT NOT NULL,
    prescribed_medications TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE,
    FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE
);

CREATE TABLE billings (
    id SERIAL PRIMARY KEY,
    appointment_id INTEGER NOT NULL UNIQUE,
    total_amount DECIMAL(10,2) NOT NULL CHECK (total_amount >= 0),
    payment_status VARCHAR(50) DEFAULT 'unpaid', 
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE
);

CREATE TABLE invoices (
    id SERIAL PRIMARY KEY,
    billing_id INTEGER NOT NULL UNIQUE,
    invoice_number VARCHAR(100) NOT NULL UNIQUE,
    consultation_fee DECIMAL(10,2) NOT NULL DEFAULT 200000,
    medication_fee DECIMAL(10,2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(10,2) NOT NULL,
    sent_to_email VARCHAR(255) NOT NULL,
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (billing_id) REFERENCES billings(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_users_email ON users (email);

CREATE INDEX idx_workouts_user_id ON workouts (user_id);

CREATE INDEX idx_exercises_workout_id ON exercises (workout_id);

CREATE INDEX idx_exercise_logs_user_id ON exercise_logs (user_id);
CREATE INDEX idx_exercise_logs_exercise_id ON exercise_logs (exercise_id);
CREATE INDEX idx_exercise_logs_created_at ON exercise_logs (created_at);

CREATE UNIQUE INDEX idx_doctors_email ON doctors (email);

CREATE INDEX idx_appointments_user_id ON appointments (user_id);
CREATE INDEX idx_appointments_doctor_id ON appointments (doctor_id);
CREATE INDEX idx_appointments_date ON appointments (appointment_date);
CREATE INDEX idx_appointments_status ON appointments (status);

CREATE INDEX idx_diagnoses_appointment_id ON diagnoses (appointment_id);
CREATE INDEX idx_diagnoses_doctor_id ON diagnoses (doctor_id);

CREATE INDEX idx_billings_appointment_id ON billings (appointment_id);
CREATE INDEX idx_billings_payment_status ON billings (payment_status);

CREATE UNIQUE INDEX idx_invoices_invoice_number ON invoices (invoice_number);
CREATE INDEX idx_invoices_billing_id ON invoices (billing_id);
CREATE INDEX idx_invoices_sent_at ON invoices (sent_at);
