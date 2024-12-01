from django.shortcuts import get_object_or_404, render
from .models import Member
from django.http import JsonResponse

def no_of_members(a, r, levels):
    if levels == 0:
        return 1
    return int(a * (r**levels - 1) / (r - 1))

def dfs_traverse_for_position(member, target_position, level, current_pos=[0]):
    if not member:
        return None

    found_member = dfs_traverse_for_position(member.left, target_position, level + 1, current_pos)
    if found_member:
        return found_member
    
    current_pos[0] += 1

    if current_pos[0] == target_position and level == 3:
        return member

    return dfs_traverse_for_position(member.right, target_position, level + 1, current_pos)

def create_tree_with_levels(request):
    if request.method == 'POST':
        root_user_id = request.POST.get('root_user_id')
        levels = int(request.POST.get('levels'))  

        root = get_object_or_404(Member, user_id=root_user_id)
        
        total_members = no_of_members(1, 2, levels)
        members = []

        for i in range(1, total_members + 1):
            user_id = f"user_{i}"
            members.append(Member.objects.create(user_id=user_id))

        def assign_member_to_tree(parent, member):
            if parent.left is None:
                parent.left = member
            elif parent.right is None:
                parent.right = member
            else:
                if parent.left:
                    assign_member_to_tree(parent.left, member)
                elif parent.right:
                    assign_member_to_tree(parent.right, member)
            parent.save()

        
        for i, member in enumerate(members):
            if i == 0:
                root.user_id = member.user_id 
                root.save() 
            else:
                assign_member_to_tree(root, member)

        return JsonResponse({
            "message": f"Tree created with {total_members} members",
            "root_user_id": root.user_id
        })
    return render(request, 'create_tree.html')








# from django.shortcuts import get_object_or_404
# from django.http import JsonResponse
# from .models import BinaryTreeMember
# from django.db import transaction

# def insert_member(request,user_id,sponsor_id=None,parent_id=None):
#     parent = None
#     position = 'root'
    
#     if parent_id:
#         parent = get_object_or_404(BinaryTreeMember,id=parent_id)
#         position = parent.available_child_position()
        
#         if not position:
#             return JsonResponse({'error':'no position available for child node'})
        
#     with transaction.atomic():
#         new_member = BinaryTreeMember.objects.create(
#             user_id = user_id,
#             position = position,
#             parent = parent,
#             sponsor_id = sponsor_id
#         )
        
#         if position == 'left':
#             parent.left_child = new_member
#         elif position == 'right':
#             parent.right_child = new_member
#         parent.save()
        
#     return JsonResponse({
#         'user_id' : new_member.userid,
#         'position' : new_member.position,
#         'parent_id' : new_member.parent_id if new_member.parent else None,
#     })
    
# def find_next_vacent_space():
#     all_members = BinaryTreeMember.objects.all()
    
#     for member in all_members:
#         if not member.has_left_child() or not member.has_right_child():
#             return member
#         return None
         
