from django.urls import path
from .views import create_tree_with_levels

urlpatterns = [
    path('', create_tree_with_levels, name='create_tree'), 
]
