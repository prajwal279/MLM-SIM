from django.contrib import admin
from .models import PackageDetails, BinaryTree, Compensations

@admin.register(PackageDetails)
class PackageDetailsAdmin(admin.ModelAdmin):
    list_display = ('base_currency', 'joining_package_fee', 'additional_product_price')

@admin.register(BinaryTree)
class BinaryTreeAdmin(admin.ModelAdmin):
    list_display = ('total_members', 'levels')

@admin.register(Compensations)
class CompensationsAdmin(admin.ModelAdmin):
    list_display = ('sponsor_bonus_percent', 'binary_pairs', 'binary_bonus_percent', 'matching_bonus_percent', 'cap_limit')
