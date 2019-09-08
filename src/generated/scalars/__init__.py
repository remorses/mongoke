
from tartiflette import Scalar
from bson import ObjectId
from typing import Union

@Scalar("Json")
class Json:
    @staticmethod
    def coerce_input(val):
        return val

    @staticmethod
    def coerce_output(val):
        return val

    def parse_literal(self, ast: "Node") -> Union[str, "UNDEFINED_VALUE"]:
        return ast.value


@Scalar("ObjectId")
class ObjectIdScalar:
    @staticmethod
    def coerce_input(val):
        return ObjectId(val)

    @staticmethod
    def coerce_output(val):
        return str(val)

    def parse_literal(self, ast: "Node") -> Union[str, "UNDEFINED_VALUE"]:
        return ast.value


@Scalar("ID")
class IDClass:
    @staticmethod
    def coerce_input(val):
        return val

    @staticmethod
    def coerce_output(val):
        return val

    def parse_literal(self, ast: "Node") -> Union[str, "UNDEFINED_VALUE"]:
        return ast.value

