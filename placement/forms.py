from django import forms

class MemberForm(forms.Form):
    num_members = forms.IntegerField(label="No Of Members",min_value=1, required=False)
    num_child = forms.IntegerField(label="No Of Children",min_value=1, required=False)
    # mat_num_child = forms.IntegerField(label="No Of Children",min_value=1, required=False)
    expense_per_user = forms.IntegerField(label="expense per user",min_value=1, required=False)
    product_name = forms.CharField(max_length=255, required=False)
    joining_package_fee = forms.CharField(max_length=255, required=False)
    b_v = forms.CharField(max_length=255, required=False, initial='0')   
    sponsor_bonus_percent = forms.FloatField(label="Sponsor Bonus (%)", min_value=0, required=False)
    # mat_sponsor_bonus_percent = forms.FloatField(label="Sponsor Bonus (%)", min_value=0, required=False)
    # uni_sponsor_bonus_percent = forms.FloatField(label="Sponsor Bonus (%)", min_value=0, required=False)
    binary_bonus_percent = forms.FloatField(label="Binary Bonus (%)", min_value=0, required=False)
    matching_bonus_percent = forms.CharField(max_length=255, required=False)
    product_quantity = forms.CharField(max_length=255, required=False)
    capping_limit = forms.FloatField(label="Capping Limit", min_value=0, required=False)
    
    BONUS_TYPE = [
        ('PRICE','PRICE'),
        ('BV','BV'),
    ]
    bonus_option = forms.ChoiceField(
    choices=BONUS_TYPE,
    label="Bonus Option",
    required=False,  
    widget=forms.Select(attrs={'class': 'form-select', 'id': 'binaryOption'})
)


    # cycle = forms.IntegerField(label="No of Cycles",min_value=1, required=False)
    CYCLE_COUNT = [
        ('S', 'weakly'),
        ('M', 'monthly'),
        ('Y', 'yearly'),
    ]
    cycle = forms.ChoiceField(choices=CYCLE_COUNT, label="Cycle", required=False, widget=forms.Select(attrs={'class': 'form-select', 'id': 'CYCLE_COUNT'}),)
    
    RATIO_CHOICES = [
        (1, '1:1'),
        (2, '1:2'),
        (3, '2:1'),
    ]
    ratio = forms.ChoiceField(choices=RATIO_CHOICES,widget=forms.RadioSelect,label="Binary Ratio", required=False)
    ratio_amount = forms.FloatField(label="Ratio Amount",min_value=1, required=False)
    BONUS_TYPE_CHOICES = [
        ('binary', 'Binary Bonus'),
        ('matching', 'Matching Bonus'),
        ('sponsor', 'Sponsor Bonus'),
    ]
    capping_scope = forms.MultipleChoiceField(choices=BONUS_TYPE_CHOICES, widget=forms.CheckboxSelectMultiple, required=False, label="Capping Scope",)
    
# class UnilevelMemberForm(forms.Form):
    # uni_num_members = forms.IntegerField(label="No Of Members",min_value=1, required=False)
    # uni_num_child = forms.IntegerField(label="No Of Children",min_value=1, required=False)
    # uni_expense_per_user = forms.IntegerField(label="expense per user",min_value=1, required=False)
    # uni_product_name = forms.CharField(max_length=255, required=False)
    # uni_joining_package_fee = forms.CharField(max_length=255, required=False)
    # uni_b_v = forms.CharField(max_length=255, required=False, initial='0')   
    # uni_sponsor_bonus_percent = forms.FloatField(label="Sponsor Bonus (%)", min_value=0, required=False)
    # uni_binary_bonus_percent = forms.FloatField(label="Binary Bonus (%)", min_value=0, required=False)
    # uni_matching_bonus_percent = forms.CharField(max_length=255, required=False)
    # uni_product_quantity = forms.CharField(max_length=255, required=False)
    # uni_capping_limit = forms.FloatField(label="Capping Limit", min_value=0, required=False)
    
#     BONUS_TYPE = [
#         ('PRICE','PRICE'),
#         ('BV','BV'),
#     ]
#     uni_bonus_option = forms.ChoiceField(
#     choices=BONUS_TYPE,
#     label="Bonus Option",
#     required=False,  
#     widget=forms.Select(attrs={'class': 'form-select', 'id': 'binaryOption'})
# )


#     # cycle = forms.IntegerField(label="No of Cycles",min_value=1, required=False)
#     CYCLE_COUNT = [
#         ('S', 'weakly'),
#         ('M', 'monthly'),
#         ('Y', 'yearly'),
#     ]
#     uni_cycle = forms.ChoiceField(choices=CYCLE_COUNT, label="Cycle", required=False, widget=forms.Select(attrs={'class': 'form-select', 'id': 'CYCLE_COUNT'}),)
    
#     RATIO_CHOICES = [
#         (1, '1:1'),
#         (2, '1:2'),
#         (3, '2:1'),
#     ]
#     uni_ratio = forms.ChoiceField(choices=RATIO_CHOICES,widget=forms.RadioSelect,label="Binary Ratio", required=False)
#     uni_ratio_amount = forms.FloatField(label="Ratio Amount",min_value=1, required=False)
#     BONUS_TYPE_CHOICES = [
#         ('binary', 'Binary Bonus'),
#         ('matching', 'Matching Bonus'),
#         ('sponsor', 'Sponsor Bonus'),
#     ]
#     uni_capping_scope = forms.MultipleChoiceField(choices=BONUS_TYPE_CHOICES, widget=forms.CheckboxSelectMultiple, required=False, label="Capping Scope",)
    
    # class MatrixMemberForm(forms.Form):
        # mat_num_members = forms.IntegerField(label="No Of Members",min_value=1, required=False)
        # mat_num_child = forms.IntegerField(label="No Of Children",min_value=1, required=False)
        # mat_expense_per_user = forms.IntegerField(label="expense per user",min_value=1, required=False)
        # mat_product_name = forms.CharField(max_length=255, required=False)
        # mat_joining_package_fee = forms.CharField(max_length=255, required=False)
        # mat_b_v = forms.CharField(max_length=255, required=False, initial='0')   
        # mat_sponsor_bonus_percent = forms.FloatField(label="Sponsor Bonus (%)", min_value=0, required=False)
    #     mat_binary_bonus_percent = forms.FloatField(label="Binary Bonus (%)", min_value=0, required=False)
    #     mat_matching_bonus_percent = forms.CharField(max_length=255, required=False)
    #     mat_product_quantity = forms.CharField(max_length=255, required=False)
    #     mat_capping_limit = forms.FloatField(label="Capping Limit", min_value=0, required=False)
        
    #     BONUS_TYPE = [
    #         ('PRICE','PRICE'),
    #         ('BV','BV'),
    #     ]
    #     mat_bonus_option = forms.ChoiceField(
    #     choices=BONUS_TYPE,
    #     label="Bonus Option",
    #     required=False,  
    #     widget=forms.Select(attrs={'class': 'form-select', 'id': 'binaryOption'})
    # )


    #     # cycle = forms.IntegerField(label="No of Cycles",min_value=1, required=False)
    #     CYCLE_COUNT = [
    #         ('S', 'weakly'),
    #         ('M', 'monthly'),
    #         ('Y', 'yearly'),
    #     ]
    #     mat_cycle = forms.ChoiceField(choices=CYCLE_COUNT, label="Cycle", required=False, widget=forms.Select(attrs={'class': 'form-select', 'id': 'CYCLE_COUNT'}),)
        
    #     RATIO_CHOICES = [
    #         (1, '1:1'),
    #         (2, '1:2'),
    #         (3, '2:1'),
    #     ]
    #     mat_ratio = forms.ChoiceField(choices=RATIO_CHOICES,widget=forms.RadioSelect,label="Binary Ratio", required=False)
    #     mat_ratio_amount = forms.FloatField(label="Ratio Amount",min_value=1, required=False)
    #     BONUS_TYPE_CHOICES = [
    #         ('binary', 'Binary Bonus'),
    #         ('matching', 'Matching Bonus'),
    #         ('sponsor', 'Sponsor Bonus'),
    #     ]
    #     mat_capping_scope = forms.MultipleChoiceField(choices=BONUS_TYPE_CHOICES, widget=forms.CheckboxSelectMultiple, required=False, label="Capping Scope",)
        
        