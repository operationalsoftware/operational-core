from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse
from fastapi.staticfiles import StaticFiles

from src.routers import contacts, home_breadcrumb
from src.ui import templates

app = FastAPI()

# routers
app.include_router(contacts.router)

# root
@app.get("/", response_class=HTMLResponse)
async def root(request: Request):
    return templates.TemplateResponse("modules/home.jinja2", {"request": request, "breadcrumbs": [home_breadcrumb] })

app.mount("/", StaticFiles(directory="static"), name="static")