export interface Stock {
  id: string
  ticker: string
  company: string
  brokerage: string
  action: string
  rating_from: string
  rating_to: string
  target_from: string
  target_to: string
  time: string
  created_at: string
  updated_at: string
}

export interface StockRecommendation {
  stock: Stock
  score: number
  reasons: string[]
  potential_up: number
}

