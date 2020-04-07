scalars_implementations = '''
from tartiflette import Scalar
from typing import Union
from tartiflette.language.ast.base import Node
from tartiflette.constants import UNDEFINED_VALUE
from tartiflette_scalars import Json, ObjectId, AnyScalar


JsonScalar = Scalar("Json")
JsonScalar(Json)

AnyScalarScalar = Scalar("AnyScalar")
AnyScalarScalar(AnyScalar)

ObjectIdScalar = Scalar("ObjectId")
ObjectIdScalar(ObjectId)


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


""" for scalar in sorted(scalars)])
}}
# print(dir(AnyScalar))
scalar_classes = [var for name, var in locals().items() if getattr(var, '_implementation', None)]
'''