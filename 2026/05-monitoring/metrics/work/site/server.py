"""Small local workshop target. No third-party Python dependencies."""
from html import escape
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
import os
from threading import Lock
from time import sleep
from urllib.parse import parse_qs, urlsplit

MODE = "normal"
LOCK = Lock()
MODES = {
    "normal": ("HTTP 200", 200, 0),
    "error": ("HTTP 500", 500, 0),
    "slow": ("Медленный ответ: 1 секунда", 200, 1),
    "timeout": ("Медленный ответ: 4 секунды", 200, 4),
}


class Handler(BaseHTTPRequestHandler):
    def send_body(self, code, body, content_type="text/html; charset=utf-8"):
        payload = body.encode("utf-8")
        try:
            self.send_response(code)
            self.send_header("Content-Type", content_type)
            self.send_header("Content-Length", str(len(payload)))
            self.send_header("Cache-Control", "no-store")
            self.end_headers()
            self.wfile.write(payload)
        except (BrokenPipeError, ConnectionResetError):
            pass  # Expected when the probe times out before the response.

    def do_GET(self):
        path = urlsplit(self.path).path
        with LOCK:
            mode = MODE
        if path == "/health":
            return self.send_body(200, "process alive\n", "text/plain; charset=utf-8")
        if path == "/site":
            label, code, delay = MODES[mode]
            sleep(delay)
            status = {
                "normal": "Сервис здоров",
                "error": "Ошибка сервиса",
                "slow": "Сервис отвечает медленно",
                "timeout": "Сервис отвечает слишком долго",
            }[mode]
            return self.send_body(code, f'<!doctype html><html lang="ru"><meta charset="utf-8"><title>Состояние сервиса</title><h1>{status}</h1><p>{escape(label)}</p><p><a href="/">Вернуться на главную</a></p></html>')
        if path != "/":
            return self.send_body(404, "not found\n", "text/plain; charset=utf-8")
        buttons = "".join(
            f'<button name="mode" value="{key}">{escape(value[0])}</button>'
            for key, value in MODES.items()
        )
        return self.send_body(200, f'''<!doctype html><html lang="ru"><meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>inview: управление сайтом</title>
<style>body{{font:18px/1.6 monospace;background:#101820;color:#eaf3f5;max-width:850px;margin:60px auto;padding:20px}}
h1{{color:#77e3c2}}button{{font:inherit;padding:12px;margin:8px;cursor:pointer}}a{{color:#ffd28a}}</style>
<h1>Восьмое чудо вайбкодинга<sup>TM</sup></h1><p>Текущий режим: <strong>{escape(MODES[mode][0])}</strong></p>
<form method="post" action="/mode">{buttons}</form>
<p><a href="/site">Открыть сайт</a></p>
</html>''')

    def do_POST(self):
        global MODE
        if urlsplit(self.path).path != "/mode":
            return self.send_body(404, "not found")
        length = int(self.headers.get("Content-Length", "0"))
        if length > 1024:
            return self.send_body(413, "too large")
        mode = parse_qs(self.rfile.read(length).decode("utf-8")).get("mode", [""])[0]
        if mode not in MODES:
            return self.send_body(400, "unknown mode")
        with LOCK:
            MODE = mode
        self.send_response(303)
        self.send_header("Location", "/")
        self.send_header("Content-Length", "0")
        self.end_headers()


if __name__ == "__main__":
    ThreadingHTTPServer(("0.0.0.0", int(os.environ.get("PORT", "8080"))), Handler).serve_forever()
