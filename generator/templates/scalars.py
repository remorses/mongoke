scalars_implementations = '''
from tartiflette import Scalar
from bson import ObjectId

@Scalar("Json")
class Json:
    @staticmethod
    def coerce_input(val):
        return val

    @staticmethod
    def coerce_output(val):
        return val

@Scalar("ObjectId")
class ObjectIdScalar:
    @staticmethod
    def coerce_input(val):
        return ObjectId(val)

    @staticmethod
    def coerce_output(val):
        return str(val)
${{
''.join([f"""
@Scalar("{scalar}")
class {scalar}Class:
    @staticmethod
    def coerce_input(val):
        return val

    @staticmethod
    def coerce_output(val):
        return val
""" for scalar in scalars])
}}
'''