from django import forms



class MLMCalculationForm(forms.Form):
    
    CURRENCY_TYPE = [
        ('INR','INR'),
        ('USD','USD'),
        ('EURO','EURO'),
        ('HKD','HKD'),
        ('AED','AED'),
    ]
    base_currency = forms.ChoiceField(label="Base Currency",choices=CURRENCY_TYPE)
    joining_package_fee = forms.FloatField(label="Joining Package Fee", min_value=0)
    additional_product_price = forms.FloatField(label="Additional Product Price", min_value=0)
    levels = forms.IntegerField(label="Number of Levels", min_value=0)
    sponsor_bonus_percent = forms.FloatField(label="Sponsor Bonus (%)", min_value=0)
    binary_pairs = forms.IntegerField(label="Number of Binary Pairs", min_value=0)
    binary_bonus_percent = forms.FloatField(label="Binary Bonus (%)", min_value=0)
    matching_bonus_percent = forms.FloatField(label="Matching Bonus (%)", min_value=0)
    matching_bonus_levels = forms.IntegerField(label="Matching Bonus Levels", min_value=0)
    cap_limit = forms.FloatField(label="Capping Limit", min_value=0)
    
    BONUS_TYPE_CHOICES = [
        ('binary', 'Binary Bonus'),
        ('matching', 'Matching Bonus'),
        ('sponsor', 'Sponsor Bonus'),
        ('total', 'Total Bonus'),
    ]
    capping_scope = forms.ChoiceField(choices=BONUS_TYPE_CHOICES)
