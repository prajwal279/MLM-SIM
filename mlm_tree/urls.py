from django.urls import path
from .views import insert_view

urlpatterns = [
    path('', insert_view, name='insert'),
]
