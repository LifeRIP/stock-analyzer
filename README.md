# 📊 Stock-Analyzer

Stock-Analyzer is an application designed to analyze stock data retrieved from an external API and provide useful insights for investing. 💼

## 🎯 Application Purpose

Provide investment recommendations based on an automated analysis of broker activity, rating changes, and price projections. 📈

## 📱 User Interface

![Main Interface](/mockups/mockup1.png)
![Recommendations and Stock Desktop Interface](/mockups/mockup2.png)
![Recommendations and Stock Mobile Interface](/mockups/mockup3.png)

## ⚙️ Requirements

- [Docker](https://www.docker.com/) installed on your system 🐳
- `.env` file in the backend directory 🛠️

## 🚀 Setup

1. **Configure the environment** 🧪:

   - In the backend directory, make sure to include a `.env` file based on the provided `.env.example`.
   - You can copy and rename it using the following command:
     ```sh
     cp backend/.env.example backend/.env
     ```
   - Then, edit the `.env` file with the necessary configurations.

2. **Start the project** 🏁:
   - From the project root, run the following command:
     ```sh
     docker compose up
     ```
   - This will start the necessary containers to run the application.

3. **Open the App Web UI** 🔄:
   - Go to [http://localhost:8082](http://localhost:8082)
   - Once there, click on the **"Sync Stocks"** button to fetch the latest stock data.


## 🖥️ Usage

Once the containers are up and running, you can access the application through the following URLs:

- **CockroachDB Web UI**: [http://localhost:8080/](http://localhost:8080/)
- **API**: [http://localhost:8081/](http://localhost:8081/)
- **App Web UI**: [http://localhost:8082/](http://localhost:8082/)

## ⚙️ How Does the Recommendation System Work?

1. **Data Collection**: It gathers all the stocks that had relevant movements on a specific date.
2. **Recommendation Processing**: Each stock is automatically evaluated based on three factors:
   - 📊 Rating changes.
   - 🧾 Actions taken by brokers (such as raising the target price).
   - 💰 Growth potential based on the target price change.
3. **Score Calculation**: Each factor is assigned a specific score. Stocks are sorted by this score and the top 5 are shown.
4. **Result**: The user receives a clear recommendation, with an explanation, score, and potential growth.

## 🧮 Algorithm Breakdown

```go
// Stock evaluation:
score = ratingScore + actionScore + potentialGrowth
```

- **RatingScore**: Assigned based on rating improvement or deterioration.
- **ActionScore**: Based on broker actions (e.g., upgraded rating, raised target).
- **PotentialGrowth**: Percentage change between the previous and new normalized target price.

## 🧾 Table of Actions and Their Impact

| 🎯 Action                | 🔢 Points | 📘 Meaning                                 |
| ------------------------ | --------- | ------------------------------------------ |
| ✅ **Upgraded by**       | +8        | Rating of the stock was improved.          |
| 🟢 **Target raised by**  | +7        | Target price was increased.                |
| 🟡 **Target set by**     | +6        | A target price was set for the first time. |
| 🟡 **Initiated by**      | +5        | Coverage of the stock was initiated.       |
| 🟡 **Reiterated by**     | +4        | Rating was reiterated.                     |
| 🔴 **Downgraded by**     | +3        | Rating was lowered.                        |
| ❌ **Target lowered by** | +1        | Target price was decreased.                |

## ⭐ Rating Table and Its Meaning

| 🏷️ Rating                 | 🔢 Points | 📘 Interpretation                           |
| ------------------------- | --------- | ------------------------------------------- |
| ✅ Strong-Buy, Buy, etc.  | +5        | High expectation of price increase.         |
| 🟢 Outperform, Overweight | +4        | Expected to perform better than the market. |
| 🟡 Hold, Neutral, etc.    | +3        | Expected to remain stable.                  |
| 🔴 Underperform, Negative | +2        | Expected to perform below the market.       |
| ❌ Sell                   | +1        | Recommendation to sell.                     |

## 💡 Example Result for a Stock

**Apple Inc. (AAPL)**

- ⭐ **Score:** 14.3
- 🔍 **Reasons:**
  - Rating improvement from Hold to Buy (+2)
  - Action taken: Target raised by JP Morgan (+7)
  - Increase in target price by 52.86% (+5.3)
- 📈 **Growth potential:** 52.86%

## ⚠️ Disclaimer

This project is intended for educational and informational purposes only.

All investment decisions are the sole responsibility of the user. I am not responsible for any financial losses or outcomes resulting from the use of this tool. Please do your own research and consult a professional financial advisor before making any investment.