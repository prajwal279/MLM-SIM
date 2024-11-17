from django import forms

class MemberForm(forms.Form):
    num_members=forms.IntegerField(label="No Of Members",min_value=1)
    joining_package_fee = forms.CharField(max_length=255, required=False)
    sponsor_bonus_percent = forms.FloatField(label="Sponsor Bonus (%)", min_value=0)
    binary_bonus_percent = forms.FloatField(label="Binary Bonus (%)", min_value=0)
    matching_bonus_percent = forms.CharField(max_length=255, required=False)
    product_quantity = forms.CharField(max_length=255, required=False)
    capping_limit = forms.FloatField(label="Capping Limit", min_value=0)
    cycle = forms.IntegerField(label="No of Cycles",min_value=1)

    BONUS_TYPE_CHOICES = [
        ('binary', 'Binary Bonus'),
        ('matching', 'Matching Bonus'),
        ('sponsor', 'Sponsor Bonus'),
        ('total', 'Total Bonus'),
    ]
    capping_scope = forms.ChoiceField(choices=BONUS_TYPE_CHOICES)
    
    CARRY_CHOICE = [
        ('yes','Yes'),
        ('no','No'),
    ]
    carry_yes_no = forms.ChoiceField(choices=CARRY_CHOICE,widget=forms.RadioSelect,label="Carry Forward (Yes or No)",required=True)
    
    # def __init__(self, *args, **kwargs):

    #     cycles = kwargs.get('cycles', 1) 
        
    #     super().__init__(*args, **kwargs)

    #     for i in range(1, cycles + 1):
    #         self.fields[f'joining_package_fee_{i}'] = forms.FloatField(
    #             label=f"Joining Package Fee for Cycle {i}", min_value=0, required=True)
