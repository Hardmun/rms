select img._fld598         as imageID,
       img._fld68          as collection,
       img._fld69          as picture,
       img._fld70          as form,
       img._fld71          as color,
       brands._code        as brandCode,
       brands._description as brand,
       img._fld64rref
from _reference17 as img
         left join _reference13 as brands
                   on img._Fld67RRef = brands._idrref
where true