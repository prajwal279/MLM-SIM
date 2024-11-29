from django.shortcuts import render
from django.views.decorators.csrf import csrf_exempt
from django.http import JsonResponse
import json
import requests
from .forms import MemberForm
 
def build_new_tree(request):
    print(request.POST)
    if request.method == 'POST':
        form = MemberForm(request.POST) 
        if form.is_valid():
            num_members = form.cleaned_data['num_members']
            expense_per_user = form.cleaned_data['expense_per_user']
            product_name = [level.strip() for level in request.POST.getlist('product_name', '') if level.strip()]
            sponsor_bonus_percent = form.cleaned_data['sponsor_bonus_percent']
            binary_bonus_percent = form.cleaned_data['binary_bonus_percent']
            joining_package_fee = [int(level.strip()) for level in request.POST.getlist('joining_package_fee', '') if level.strip().isdigit()]
            b_volume = [int(level.strip()) for level in request.POST.getlist('b_v', '') if level.strip().isdigit()]
            bonus_option = form.cleaned_data['bonus_option']
            product_quantity = [int(level.strip()) for level in request.POST.getlist('product_quantity', '') if level.strip().isdigit()]
            capping_limit = form.cleaned_data['capping_limit']
            capping_scope = ','.join(form.cleaned_data['capping_scope'])
            matching_bonus_percents = [int(level.strip()) for level in request.POST.getlist('matching_bonus_percent', '') if level.strip().isdigit()]
            cycle = form.cleaned_data['cycle']
            ratio = form.cleaned_data['ratio']
            ratio_amount = form.cleaned_data['ratio_amount']
            data = {
                "num_members": num_members,
                "expense_per_user": expense_per_user,
                "product_name": product_name,                                            
                "sponsor_percentage": sponsor_bonus_percent,
                "binary_percentage": binary_bonus_percent,
                "joining_package_fee": joining_package_fee,
                "b_volume" : b_volume,
                "bonus_option": bonus_option,
                "product_quantity": product_quantity,
                "capping_amount": capping_limit,
                "capping_scope": capping_scope,
                "matching_percentage": matching_bonus_percents,
                "cycle": cycle,
                "ratio":ratio,
                "ratio_amount":ratio_amount,
            }
           
            try:
                response = requests.post('http://localhost:9000/calculate', json=data)
                response.raise_for_status() 

                results = response.json()  
                return render(request, 'display_members.html', {
                    'results': results,
                    'num_members':num_members,
                    'joining_package_fee':joining_package_fee,
                })
            except requests.exceptions.RequestException as e:
                return JsonResponse({'error': f'Failed to communicate with Go server: {str(e)}'}, status=500)

        else:
            return render(request, 'interface.html', {'form': form})
    else:
        form = MemberForm()  
        return render(request, 'interface.html', {'form': form})





def build_unilevel_tree(request):
    print(request.POST)
    if request.method == 'POST':
        form = MemberForm(request.POST) 
        if form.is_valid():
            num_members = form.cleaned_data['num_members']
            num_child = form.cleaned_data['num_child']
            expense_per_user = form.cleaned_data['expense_per_user']
            product_name = [level.strip() for level in request.POST.getlist('product_name', '') if level.strip()]
            sponsor_bonus_percent = form.cleaned_data['sponsor_bonus_percent']
            joining_package_fee = [int(level.strip()) for level in request.POST.getlist('joining_package_fee', '') if level.strip().isdigit()]
            bonus_option = form.cleaned_data['bonus_option']
            capping_limit = form.cleaned_data['capping_limit']
            capping_scope = ','.join(form.cleaned_data['capping_scope'])
            matching_bonus_percents = [int(level.strip()) for level in request.POST.getlist('matching_bonus_percent', '') if level.strip().isdigit()]
            cycle = form.cleaned_data['cycle']
            product_quantity = [int(level.strip()) for level in request.POST.getlist('product_quantity', '') if level.strip().isdigit()]
            data = {
                "num_members": num_members,
                "num_child": num_child,
                "expense_per_user": expense_per_user,
                "product_name": product_name,                                            
                "sponsor_percentage": sponsor_bonus_percent,
                "joining_package_fee": joining_package_fee,
                "bonus_option": bonus_option,
                "product_quantity": product_quantity,
                "capping_amount": capping_limit,
                "capping_scope": capping_scope,
                "matching_percentage": matching_bonus_percents,
                "cycle": cycle,
            }
            try:
                response = requests.post('http://localhost:9000/unilevel', json=data)
                response.raise_for_status() 

                results = response.json()  
                return render(request, 'display_members.html', {
                    'results': results,
                    'num_members':num_members,
                })
            except requests.exceptions.RequestException as e:
                return JsonResponse({'error': f'Failed to communicate with Go server: {str(e)}'}, status=500)

        else:
            return render(request, 'interface.html', {'form': form})
    else:
        form = MemberForm()  
        return render(request, 'interface.html', {'form': form})



@csrf_exempt
def process_results(request):
    if request.method == 'POST':
        try:
            data = json.loads(request.body)
            
        except json.JSONDecodeError:
            return JsonResponse({'error': 'Invalid JSON data'}, status=400)
        cycles = data.get('cycles', {})
        nodes = data.get('tree_structure', [])
        context={}
        context['sponsor_bonus']="sponsor_bonus"
        context['binary_bonus']="binary_bonus"
        context['nodes']=nodes
        context['cycles']=cycles
        
        render_context = render(request, 'display_members.html', context)
        return render_context
    else:
        return JsonResponse({'error': 'Invalid request method'}, status=405)
    
    
    