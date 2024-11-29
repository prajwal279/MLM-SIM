# from django.urls import path
# from .views import build_new_tree

# urlpatterns = [
#     path('', build_new_tree, name='build_new_tree'),
# ]
from django.urls import path
from .views import build_new_tree, process_results, build_unilevel_tree  # Import your views

urlpatterns = [
    path('', build_new_tree, name='build_new_tree'),
    path('unilevel/', build_unilevel_tree, name='build_unilevel_tree'),
    path('process_results/', process_results, name='process_results'),
]
