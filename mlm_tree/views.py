from django.http import JsonResponse
from collections import deque

class Node:
    def __init__(self, x):
        self.data = x
        self.left = None
        self.right = None

def InsertNode(root, data):
    if root is None:
        return Node(data)

    q = deque()
    q.append(root)

    while q:
        curr = q.popleft()

        if curr.left is not None:
            q.append(curr.left)
        else:
            curr.left = Node(data)
            return root

        if curr.right is not None:
            q.append(curr.right)
        else:
            curr.right = Node(data)
            return root


def inorder(curr):
    if curr is None:
        return []
    return inorder(curr.left) + [curr.data] + inorder(curr.right)

def preorder(curr):
    if curr is None:
        return []
    return [curr.data] + preorder(curr.left) + preorder(curr.right)

def postorder(curr):
    if curr is None:
        return []
    return  postorder(curr.right) + postorder(curr.left) + [curr.data]

def insert_view(request):
    key = int(request.GET.get('key', 12)) 

    root = Node(10)
    root.left = Node(11)
    root.right = Node(9)
    root.left.left = Node(7)
    root.right.left = Node(15)
    root.right.right = Node(8)

    root = InsertNode(root, key)

    result = inorder(root)
    r = preorder(root)
    res = postorder(root)
    return JsonResponse({'inorder': result,
                         'preorder' : r,
                         'postorder':res})
