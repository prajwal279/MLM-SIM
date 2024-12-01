from django.contrib import admin
from .models import Member

@admin.register(Member)
class MemberAdmin(admin.ModelAdmin):
    list_display = ('user_id', 'position', 'parent', 'left', 'right', 'sponsor')  
    search_fields = ('user_id', 'parent__user_id')  

