from django.shortcuts import render
from .models import PackageDetails, BinaryTree, Compensations
from .forms import MLMCalculationForm

def no_of_member(a,r,levels):
    if  levels == 0:
        return 1
    else:
        return a*(r**levels-1)/(r-1)

def calculate_mlm_logic(package_details, binary_tree, compensations, bonus_type):
    income_from_joining = binary_tree.total_members * package_details.joining_package_fee
    income_from_product = binary_tree.total_members * package_details.additional_product_price
    
    sponsor_bonus = round(binary_tree.total_members * package_details.joining_package_fee * (compensations.sponsor_bonus_percent / 100),2)
    
    binary_bonus = round(package_details.joining_package_fee * (compensations.binary_bonus_percent / 100),2)
    matching_bonus = 0
    if compensations.matching_bonus_levels <= binary_tree.levels:
        for i in range(compensations.matching_bonus_levels):
            matching_bonus = round(binary_bonus * (compensations.matching_bonus_percent / 100)*i,2) 
    else:
        matching_bonus = 0
    total_bonus = round(sponsor_bonus + binary_bonus + matching_bonus,2)
   
    if bonus_type == 'sponsor':
        total_bonus = round(sponsor_bonus,2)
    elif bonus_type == 'binary':
        total_bonus = round(binary_bonus,2)
    elif bonus_type == 'matching':
        total_bonus = round(matching_bonus,2)
    else:
        total_bonus = round(sponsor_bonus + binary_bonus + matching_bonus,2)
    capped_bonus = round(max(total_bonus-compensations.cap_limit,0),2)
    company_profit = income_from_joining + income_from_product - capped_bonus

    return {
        'income_from_joining': income_from_joining,
        'income_from_product': income_from_product,
        'sponsor_bonus': sponsor_bonus,
        'binary_bonus': binary_bonus,
        'matching_bonus': matching_bonus,
        'total_bonus': total_bonus,
        'bonus_type': bonus_type,
        'capped_bonus': capped_bonus,
        'company_profit': company_profit,
    }

def mlm_calculate_view(request):
    if request.method == 'POST':
        form = MLMCalculationForm(request.POST)
        
        if form.is_valid():
            base_currency = form.cleaned_data['base_currency']
            joining_package_fee = form.cleaned_data['joining_package_fee']
            additional_product_price = form.cleaned_data['additional_product_price']
            levels = form.cleaned_data['levels']
            sponsor_bonus_percent = form.cleaned_data['sponsor_bonus_percent']
            binary_pairs = form.cleaned_data['binary_pairs']
            binary_bonus_percent = form.cleaned_data['binary_bonus_percent']
            matching_bonus_percent = form.cleaned_data['matching_bonus_percent']
            matching_bonus_levels = form.cleaned_data['matching_bonus_levels']
            bonus_type=form.cleaned_data['bonus_type']
            cap_limit = form.cleaned_data['cap_limit']
            total_members = no_of_member(2, 2, levels)
            
            package_details = PackageDetails.objects.create(
                base_currency=base_currency,
                joining_package_fee=joining_package_fee,
                additional_product_price=additional_product_price
            )
            binary_tree = BinaryTree.objects.create(
                levels=levels,
                total_members=total_members
            )   
            compensations = Compensations.objects.create(
                sponsor_bonus_percent=sponsor_bonus_percent,
                binary_pairs=binary_pairs,
                binary_bonus_percent=binary_bonus_percent,
                matching_bonus_percent=matching_bonus_percent,
                matching_bonus_levels=matching_bonus_levels,
                cap_limit=cap_limit
            )
            result = calculate_mlm_logic(package_details, binary_tree, compensations, bonus_type)
            return render(request, 'result.html', {'result': result , 'total_members' : total_members})
        else:
            return render(request, 'mlm.html', {'form': form})
    else:
        form = MLMCalculationForm()
        return render(request, 'mlm.html', {'form': form})
    