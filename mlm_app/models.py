from django.db import models


class PackageDetails(models.Model):
    base_currency = models.CharField(max_length=10)
    joining_package_fee = models.FloatField()
    additional_product_price = models.FloatField()

    def __str__(self):
        return f"Base Currency: {self.base_currency}, Joining Fee: {self.joining_package_fee}, Product Price: {self.additional_product_price}"

class BinaryTree(models.Model):
    total_members = models.IntegerField()
    levels = models.IntegerField()

    def __str__(self):
        return f"Levels: {self.levels}, Total Members: {self.total_members}"

class Compensations(models.Model):
    sponsor_bonus_percent = models.FloatField()
    bonus_type = models.FloatField(default=0.0)
    binary_pairs = models.IntegerField()
    binary_bonus_percent = models.FloatField()
    matching_bonus_percent = models.FloatField()
    matching_bonus_levels = models.IntegerField()
    
    cap_limit = models.FloatField()

    def __str__(self):
        return f"Sponsor Bonus: {self.sponsor_bonus_percent}%, Binary Bonus: {self.binary_bonus_percent}%, Matching Bonus: {self.matching_bonus_percent}%"
