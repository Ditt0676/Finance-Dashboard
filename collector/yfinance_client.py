"""
Cliente para Yahoo Finance via yfinance
"""

import logging
import pandas as pd
import yfinance as yf
from datetime import datetime

logger = logging.getLogger(__name__)


def fetch_stocks(symbols: list[str]) -> dict:
    result = {}

    for symbol in symbols:
        try:
            ticker = yf.Ticker(symbol)

            info = ticker.info
            current_price = info.get("currentPrice")
            if current_price is None:
                current_price = info.get("regularMarketPrice")

            hist = ticker.history(period="30d")
            history = []
            for idx, row in hist.iterrows():
                history.append({
                    "date": idx.strftime("%Y-%m-%d"),
                    "open": round(row["Open"], 2),
                    "high": round(row["High"], 2),
                    "low": round(row["Low"], 2),
                    "close": round(row["Close"], 2),
                    "volume": int(row["Volume"]),
                })
            
            result[symbol] = {
                "symbol": symbol,
                "name": info.get("longName", symbol),
                "current_price": current_price,
                "currency": info.get("currency", "USD"),
                "market_cap": info.get("marketCap"),
                "52w_high": info.get("fiftyTwoWeekHigh"),
                "52w_low": info.get("fiftyTwoWeekLow"),
                "history": history,
                "updated_at": datetime.now().isoformat(),
            }

        except Exception as e:
            logger.error("Error fetching %s: %s", symbol, e)
            result[symbol] = {"symbol": symbol, "error": str(e)}

    return result
