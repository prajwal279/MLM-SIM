from django.db import models

class Member(models.Model):
    user_id = models.CharField(max_length=100, unique=True)
    position = models.IntegerField() 
    parent = models.ForeignKey('self', on_delete=models.SET_NULL, related_name='children', null=True, blank=True)
    left = models.OneToOneField('self', on_delete=models.SET_NULL, related_name='left_child', null=True, blank=True)
    right = models.OneToOneField('self', on_delete=models.SET_NULL, related_name='right_child', null=True, blank=True)
    sponsor = models.ForeignKey('self', on_delete=models.SET_NULL, related_name='sponsored_members', null=True, blank=True)

    def __str__(self):
        return f"Member {self.user_id} (Position: {self.position})"

    def has_vacancy(self):
        return self.left is None or self.right is None
