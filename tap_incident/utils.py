from datetime import datetime

def safe_fromisoformat(date_str) -> datetime:
    """
    Convert a date string to a datetime object.
    
    Convert 2025-05-07T17:30:55.97Z -> 2025-05-07T17:30:55.000000+00:00
    Incident.io can return dates with the wrong number of decimal places for fromisoformat. 
    Below python 3.11 this will fail as an argument to datetime.fromisoformat()
    """
    if "." in date_str:
        head, _ = date_str.split(".", 1)
        zulu_date_with_zeroed_ms = head + ".000000+00:00" # ensure UTC with 6 decimal places
        return datetime.fromisoformat(zulu_date_with_zeroed_ms)
    else:
        return datetime.fromisoformat(date_str)