from fastapi import APIRouter, Request
from fastapi.responses import HTMLResponse

from ..ui import templates
from . import home_breadcrumb

router = APIRouter(prefix="/contacts")

CONTACTS = [
    {"id": 1, "name": "John Doe", "email": "jon.doe@example.com"},
    {"id": 2, "name": "Jane Doe", "email": "jane.doe@example.com"},
]

contacts_breadcrumb = {"name": "Contacts", "url": "/contacts", "icon": "account-group"}


@router.get("", response_class=HTMLResponse)
async def contact_root(request: Request):
    return templates.TemplateResponse("modules/contacts/index.jinja2", {
        "request": request,
        "breadcrumbs": [
            home_breadcrumb,
            contacts_breadcrumb
            ],
        "contacts": CONTACTS
        })


@router.get("/{id}", response_class=HTMLResponse)
async def contact_show(request: Request, id: int):
    return templates.TemplateResponse("modules/contacts/show.jinja2", {
        "request": request,
        "breadcrumbs": [
            home_breadcrumb,
            contacts_breadcrumb,
            {
                "name": CONTACTS[id - 1]["name"],
                "url": f"/contacts/{id}",
                "icon": "account"
            }],
        "contact": CONTACTS[id - 1],
        })
