"""
Cliente para CoinGecko API
"""

import logging
import requests
from datetime import datetime
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry

logger = logging.getLogger(__name__)

BASE_URL = "https://api.coingecko.com/api/v3"

_session = None


def _get_session():
    global _session
    if _session is None:
        _session = requests.Session()
        retries = Retry(
            total=3,
            backoff_factor=0.5,
            status_forcelist=[429, 500, 502, 503, 504],
            allowed_methods=["GET"],
        )
        _session.mount("https://", HTTPAdapter(max_retries=retries))
    return _session


def fetch_crypto(coin_ids: list[str]) -> dict:
    result = {}
    session = _get_session()

    try:
        ids_param = ",".join(coin_ids)
        response = session.get(
            f"{BASE_URL}/simple/price",
            params={
                "ids": ids_param,
                "vs_currencies": "usd",
                "include_market_cap": "true",
                "include_24hr_change": "true",
                "include_24hr_vol": "true",
            },
            timeout=10,
        )
        response.raise_for_status()
        prices = response.json()

    except Exception as e:
        logger.error("Error fetching crypto prices: %s", e)
        prices = {}

    for coin_id in coin_ids:
        try:
            params = {"vs_currency": "usd", "days": "30"}
            hist_response = session.get(
                f"{BASE_URL}/coins/{coin_id}/market_chart",
                params=params,
                timeout=10,
            )
            hist_response.raise_for_status()
            hist_data = hist_response.json()

            seen = set() 

            history = []
            
            for p in hist_data.get ("prices", []):
                date = datetime.fromtimestamp(p[0] / 1000).strftime("%Y-%m-%d")
                if date not in seen:
                    seen.add(date)
                    history.append({"date": date, "price": round(p [1], 4)})

            price_info = prices.get(coin_id, {})

            result[coin_id] = {
                "id": coin_id,
                "current_price": price_info.get("usd"),
                "market_cap": price_info.get("usd_market_cap"),
                "volume_24h": price_info.get("usd_24h_vol"),
                "change_24h": price_info.get("usd_24h_change"),
                "history": history,
                "updated_at": datetime.utcnow().isoformat(),
            }

        except Exception as e:
            logger.error("Error fetching %s: %s", coin_id, e)
            result[coin_id] = {"id": coin_id, "error": str(e)}

    return result
