"""
Utilidad para guardar datos como JSON local
"""

import json
import logging
from datetime import datetime
from pathlib import Path


def save_json(data: dict, filepath: str) -> None:
    """
    Guarda data como JSON. Crea el directorio si no existe.
    También guarda un backup con timestamp cada vez.
    """
    path = Path(filepath)
    path.parent.mkdir(parents=True, exist_ok=True)

    # Archivo principal (siempre sobreescrito)
    with open(path, "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2, ensure_ascii=False)

    # Backup con timestamp (para histórico)
    backup_dir = path.parent / "backups"
    backup_dir.mkdir(exist_ok=True)
    timestamp = datetime.now().strftime("%Y%m%d_%H%M")
    backup_path = backup_dir / f"{path.stem}_{timestamp}.json"
    with open(backup_path, "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2, ensure_ascii=False)

    # Mantener solo los últimos 10 backups por archivo
    backups = sorted(backup_dir.glob(f"{path.stem}_*.json"))
    for old in backups[:-10]:
        old.unlink()


def load_json(filepath: str) -> dict:
    """Carga un JSON local. Retorna dict vacío si no existe."""
    path = Path(filepath)
    if not path.exists():
        return {}
    try:
        with open(path, "r", encoding="utf-8") as f:
            return json.load(f)
    except json.JSONDecodeError as e:
        logging.warning("JSON corrupto en %s: %s", filepath, e)
        return {}