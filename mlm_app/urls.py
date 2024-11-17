from django.contrib import admin
from django.urls import path
from . import views

urlpatterns = [
    path('',views.mlm_calculate_view,name="mlm_calculate_view"),
]


