from django import forms 

class Binary_tree(forms.Form):
    levels = forms.IntegerField(label="No Of Levels",max_length=20)
    user_id = forms.CharField(max_length=100, unique=True)
    position = forms.IntegerField() 
    parent = forms.ForeignKey('self', on_delete=forms.SET_NULL, related_name='children', null=True, blank=True)
    left = forms.OneToOneField('self', on_delete=forms.SET_NULL, related_name='left_child', null=True, blank=True)
    right = forms.OneToOneField('self', on_delete=forms.SET_NULL, related_name='right_child', null=True, blank=True)
    sponsor = forms.ForeignKey('self', on_delete=forms.SET_NULL, related_name='sponsored_members', null=True, blank=True)