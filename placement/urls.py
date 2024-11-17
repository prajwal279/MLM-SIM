# from django.urls import path
# from .views import build_new_tree

# urlpatterns = [
#     path('', build_new_tree, name='build_new_tree'),
# ]
from django.urls import path
from .views import build_new_tree, process_results  # Import your views

urlpatterns = [
    path('', build_new_tree, name='build_new_tree'),
    path('process_results/', process_results, name='process_results'),  # Add this line
]
