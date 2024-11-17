# from django.test import TestCase
# from django.urls import reverse
# from .models import Tree_structure
# from .forms import MemberForm
# from decimal import Decimal

# class TreeStructureTests(TestCase):

#     def setUp(self):
#         self.test_user_1 = Tree_structure.objects.create(userid=1, parentid=None, position=None, levels=0, lft=1, rgt=2)
#         self.test_user_2 = Tree_structure.objects.create(userid=2, parentid=self.test_user_1, position="left", levels=1, lft=3, rgt=4)
#         self.test_user_3 = Tree_structure.objects.create(userid=3, parentid=self.test_user_1, position="right", levels=1, lft=5, rgt=6)
#         self.test_user_4 = Tree_structure.objects.create(userid=4, parentid=self.test_user_2, position="left", levels=2, lft=7, rgt=8)
#         self.test_user_5 = Tree_structure.objects.create(userid=5, parentid=self.test_user_2, position="right", levels=2, lft=9, rgt=10)

#     def test_create_new_tree(self):
#         data = {
#             'num_members': 5,
#             'joining_package_fee': 1000,
#             'sponsor_bonus_percent': 10,
#             'binary_bonus_percent': 5,
#             'capping_limit': 5000,
#             'capping_scope': 'total',
#             'carry_yes_no': 'no',
#             'matching_bonus_percent': '5,10,15',
#         }

#         form = MemberForm(data)
#         self.assertTrue(form.is_valid())
        
#         response = self.client.post(reverse('build_new_tree'), data)
#         self.assertEqual(response.status_code, 200)
#         self.assertTemplateUsed(response, 'display_members.html')
#         self.assertEqual(Tree_structure.objects.count(), 5)

#     def test_calculate_sponsor_bonus(self):
#         nodes = Tree_structure.objects.all()
#         sponsor_bonus_percent = 10
#         joining_package_fee = 1000
#         expected_bonus = 2 * sponsor_bonus_percent / 100 * joining_package_fee 
        
#         total_bonus = self.test_calculate_sponsor_bonus(nodes, sponsor_bonus_percent, joining_package_fee)
        
#         self.assertEqual(total_bonus, expected_bonus)
#         self.assertEqual(self.test_user_1.sponsor_bonus, expected_bonus)
#         self.assertEqual(self.test_user_2.sponsor_bonus, sponsor_bonus_percent / 100 * joining_package_fee)

#     def test_calculate_binary_bonus(self):
#         """Test binary bonus calculation"""
#         nodes = Tree_structure.objects.all()
#         binary_bonus_percent = 5
#         joining_package_fee = 1000

#         # Left and right counts under the root user_1
#         expected_binary_bonus = 2 * joining_package_fee * binary_bonus_percent / 100  # Total for both children
        
#         total_bonus = self.test_calculate_binary_bonus(nodes, binary_bonus_percent, joining_package_fee)
        
#         self.assertEqual(total_bonus, expected_binary_bonus)
#         self.assertEqual(self.test_user_1.binary_bonus, expected_binary_bonus)

#     def test_calculate_matching_bonus(self):
#         """Test matching bonus calculation"""
#         nodes = Tree_structure.objects.all()
#         matching_bonus_percent = {1: 5, 2: 10}  # Level 1 gets 5%, Level 2 gets 10%
        
#         # Assuming that the binary_bonus is already calculated
#         # The matching bonus will be a function of the binary_bonus of the children
#         total_bonus = self.test_calculate_matching_bonus(nodes, matching_bonus_percent)
        
#         # Check that the total matching bonus is correctly calculated
#         # (for this test, it's expected that the logic needs to be more complex and data-driven)
#         self.assertIsInstance(total_bonus, float)

#     def test_carry_forward(self):
#         """Test carry forward logic for binary bonus"""
#         nodes = Tree_structure.objects.all()
#         binary_bonus_percent = 5
#         joining_package_fee = 1000
#         carry_yes_no = 'yes'
#         capping_limit = 5000
        
#         # First calculate binary bonuses
#         self.test_calculate_binary_bonus(nodes, binary_bonus_percent, joining_package_fee, capping_limit)
        
#         # Now calculate carry forward
#         nodes_with_carry = self.test_calculate_carry_forward(nodes, binary_bonus_percent, joining_package_fee, carry_yes_no, capping_limit)
        
#         # Check that the carry value is applied correctly
#         for node_left, node_right in nodes_with_carry:
#             if node_left and node_left.carry > 0:
#                 self.assertGreater(node_left.carry, 0)
#             if node_right and node_right.carry > 0:
#                 self.assertGreater(node_right.carry, 0)

#     def test_bonus_calculations_after_tree_build(self):
#         """Test if all bonuses are correctly calculated after building the tree"""
#         data = {
#             'num_members': 5,
#             'joining_package_fee': 1000,
#             'sponsor_bonus_percent': 10,
#             'binary_bonus_percent': 5,
#             'capping_limit': 5000,
#             'capping_scope': 'total',
#             'carry_yes_no': 'no',
#             'matching_bonus_percent': '5,10,15',
#         }
        
#         # Simulate form submission
#         response = self.client.post(reverse('build_new_tree'), data)
        
#         # Test that the total bonuses are calculated correctly
#         nodes = Tree_structure.objects.all()
#         sponsor_bonus = self.test_calculate_sponsor_bonus(nodes, 10, 1000)
#         binary_bonus = self.test_calculate_binary_bonus(nodes, 5, 1000)
#         matching_bonus = self.test_calculate_matching_bonus(nodes, {1: 5, 2: 10})

#         # Check if bonuses were correctly applied
#         self.assertGreater(sponsor_bonus, 0)
#         self.assertGreater(binary_bonus, 0)
#         self.assertGreater(matching_bonus, 0)



from django.test import TestCase
from django.urls import reverse
from .models import Tree_structure
from .forms import MemberForm
from .views import (
    calculate_sponsor_bonus,
    calculate_binary_bonus,
    calculate_matching_bonus,
    calculate_carry_forward,
    add_node
)

class TreeStructureTests(TestCase):
    def setUp(self):
        self.root_node = Tree_structure.objects.create(userid=1, levels=0, lft=1, rgt=2)
        self.node2 = Tree_structure.objects.create(userid=2, parentid=self.root_node, position="left", levels=1, lft=2, rgt=3)
        self.node3 = Tree_structure.objects.create(userid=3, parentid=self.root_node, position="right", levels=1, lft=4, rgt=5)

    def test_calculate_sponsor_bonus(self):
        nodes = Tree_structure.objects.all()
        sponsor_bonus_percent = 10
        joining_package_fee = 1000
        capping_limit = 5000

        sponsor_bonus = calculate_sponsor_bonus(nodes, sponsor_bonus_percent, joining_package_fee, capping_limit)
        self.assertIsNotNone(sponsor_bonus)
        for node in nodes:
            self.assertLessEqual(node.sponsor_bonus, capping_limit)

    def test_calculate_binary_bonus(self):
        nodes = Tree_structure.objects.all()
        joining_package_fee = 1000
        binary_bonus_percent = 10
        capping_limit = 5000

        binary_bonus = calculate_binary_bonus(nodes, joining_package_fee, binary_bonus_percent, capping_limit)
        self.assertIsNotNone(binary_bonus)
        for node in nodes:
            self.assertLessEqual(node.binary_bonus, capping_limit)

    def test_calculate_matching_bonus(self):
        nodes = Tree_structure.objects.all()
        matching_bonus_percent = {1: 10, 2: 5}
        capping_limit = 5000

        matching_bonus = calculate_matching_bonus(nodes, matching_bonus_percent, capping_limit)
        self.assertIsNotNone(matching_bonus)
        for node in nodes:
            self.assertLessEqual(node.matching_bonus, capping_limit)

    def test_carry_forward(self):
        nodes = Tree_structure.objects.all()
        binary_bonus_percent = 5
        joining_package_fee = 1000
        carry_yes_no = 'yes'
        capping_limit = 500
        
        calculate_binary_bonus(nodes, binary_bonus_percent, joining_package_fee, capping_limit)
        
        nodes_with_carry = calculate_carry_forward(nodes, binary_bonus_percent, joining_package_fee, carry_yes_no, capping_limit)
        
        for node_left, node_right in nodes_with_carry:
            if node_left and node_left.carry > 0:
                self.assertGreater(node_left.carry, 0)
            if node_right and node_right.carry > 0:
                self.assertGreater(node_right.carry, 0)


    def test_add_node(self):
        initial_count = Tree_structure.objects.count()
        add_node(4)
        self.assertEqual(Tree_structure.objects.count(), initial_count + 1)

    def test_build_new_tree_view(self):
        response = self.client.post(reverse('build_new_tree'), {
            'num_members': 3,
            'joining_package_fee': 1000,
            'sponsor_bonus_percent': 10,
            'binary_bonus_percent': 10,
            'capping_limit': 5000,
            'capping_scope': 'total',
            'matching_bonus_percent': "10,5"
        })
        self.assertEqual(response.status_code, 200)
        self.assertTemplateUsed(response, 'input.html')
        nodes = Tree_structure.objects.all()
        self.assertGreater(len(nodes), 0)
