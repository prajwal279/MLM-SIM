from django.contrib import admin
from django.urls import path, include
from mlm_app.views import mlm_calculate_view
from binary_tree.views import create_tree_with_levels
urlpatterns = [
    path('admin/', admin.site.urls),
    path('xyz', include('mlm_app.urls')),
    path('sample/', include('binary_tree.urls')),
    path('tree/', include('mlm_tree.urls')),
    path('', include('placement.urls')),
]
