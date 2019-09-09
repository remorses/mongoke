scalars_implementations = '''
from tartiflette import Scalar
from bson import ObjectId
from typing import Union

JsonScalar = Scalar("Json")
@JsonScalar
class JsonClass:
    @staticmethod
    def coerce_input(val):
        return val

    @staticmethod
    def coerce_output(val):
        return val

    def parse_literal(self, ast: "Node") -> Union[str, "UNDEFINED_VALUE"]:
        return self.coerce_input(ast.value)

AnyScalarScalar = Scalar("AnyScalar")
@AnyScalarScalar
class AnyScalarClass:
    @staticmethod
    def coerce_input(val):
        if val == 'true':
            return True
        elif val == 'false':
            return False
        else:
            try:
                return float(val)
            except Exception:
                return str(val)

    @staticmethod
    def coerce_output(val):
        return val

    def parse_literal(self, ast: "Node") -> Union[str, "UNDEFINED_VALUE"]:
        return self.coerce_input(ast.value)

ObjectIdScalar = Scalar("ObjectId")
@ObjectIdScalar
class ObjectIdClass:
    @staticmethod
    def coerce_input(val):
        return ObjectId(val)

    @staticmethod
    def coerce_output(val):
        return str(val)

    def parse_literal(self, ast: "Node") -> Union[str, "UNDEFINED_VALUE"]:
        return self.coerce_input(ast.value)


${{
''.join([f"""
{scalar}Scalar = Scalar("{scalar}")
@{scalar}Scalar
class {scalar}Class:
    @staticmethod
    def coerce_input(val):
        return val

    @staticmethod
    def coerce_output(val):
        return val

    def parse_literal(self, ast: "Node") -> Union[str, "UNDEFINED_VALUE"]:
        return self.coerce_input(ast.value)


""" for scalar in scalars])
}}
# print(dir(AnyScalar))
scalar_classes = [var for name, var in locals().items() if getattr(var, '_implementation', None)]
'''