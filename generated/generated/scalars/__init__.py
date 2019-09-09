
from tartiflette import Scalar
from bson import ObjectId
from typing import Union
from tartiflette.language.ast.base import Node
from tartiflette.constants import UNDEFINED_VALUE


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



IDScalar = Scalar("ID")
@IDScalar
class IDClass:
    @staticmethod
    def coerce_input(val):
        return val

    @staticmethod
    def coerce_output(val):
        return val

    def parse_literal(self, ast: "Node") -> Union[str, "UNDEFINED_VALUE"]:
        return self.coerce_input(ast.value)



# print(dir(AnyScalar))
scalar_classes = [var for name, var in locals().items() if getattr(var, '_implementation', None)]
