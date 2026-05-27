"""
Finance Dashboard - Collector
Fase 1: Recolecta y guarda datos financieros localmente
"""

import os
import time
import logging
import schedule
from dotenv import load_dotenv
from yfinance_client import fetch_stocks
from coingecko_client import fetch_crypto
from save_data import save_json

load_dotenv()

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)

STOCKS =[s.strip() for s in os.getenv("STOCKS", "AAPL,GOOGL,MSFT,TSLA").split(",")]
CRYPTOS =[s.strip() for s in os.getenv("CRYPTOS", "bitcoin,ethereum,solana").split(",")]


def collect_all():
    logging.info("Recolectando datos...")

    stocks_data = fetch_stocks(STOCKS)
    save_json(stocks_data, "data/stocks.json")
    logging.info("Acciones guardadas: %s", list(stocks_data.keys()))

    crypto_data = fetch_crypto(CRYPTOS)
    save_json(crypto_data, "data/crypto.json")
    logging.info("Crypto guardada: %s", list(crypto_data.keys()))

    logging.info("Listo.\n")


if __name__ == "__main__":
    collect_all()

    schedule.every(60).minutes.do(collect_all)

    logging.info("Scheduler activo. Actualizando cada 60 minutos...")
    logging.info("Presiona Ctrl+C para detener.\n")

    while True:
        schedule.run_pending()
        time.sleep(1)
