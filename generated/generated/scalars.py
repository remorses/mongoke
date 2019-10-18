
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



# print(dir(AnyScalar))
scalar_classes = [var for name, var in locals().items() if getattr(var, '_implementation', None)]
