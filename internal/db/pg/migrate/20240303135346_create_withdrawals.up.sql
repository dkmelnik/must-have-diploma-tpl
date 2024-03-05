CREATE TABLE IF NOT EXISTS withdrawals (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  order_number VARCHAR(255) NOT NULL,
  amount DECIMAL(8,2) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);