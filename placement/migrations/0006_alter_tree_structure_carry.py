# Generated by Django 5.1.2 on 2024-11-08 05:33

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('placement', '0005_alter_tree_structure_carry'),
    ]

    operations = [
        migrations.AlterField(
            model_name='tree_structure',
            name='carry',
            field=models.IntegerField(default=0),
        ),
    ]