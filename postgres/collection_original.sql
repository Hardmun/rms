with data_res as (select products._Fld36          as collection,
                         encode(products._IDRRef, 'hex') as uuid,
                         products._Code                as code,
                         products._Description         as description,
                         products._Fld37              as length,
                         products._Fld38               as width,
                         -sum(reserved._Fld243)      as count
                  from _AccumRgT244 as reserved
                           inner join _Reference12 as products
                                      on reserved._Fld239RRef = products._IDRRef
                  where true
                  group by products._Fld36, products._IDRRef, products._Code, products._Description, products._Fld37,
                           products._Fld38, reserved._Fld240RRef
                  having sum(reserved._Fld243) > 0),

     data_collection as
         (select products._Fld36          as collection,
                 encode(products._IDRRef, 'hex') as uuid,
                 products._Code                as code,
                 products._Description         as description,
                 products._Fld37              as length,
                 products._Fld38               as width,
                 sum(balance._Fld249)       as count
          from _AccumRgT251 as balance
                   inner join _Reference12 as products
                              on balance._Fld246RRef = products._IDRRef
          where true
          group by products._Fld36, products._IDRRef, products._Code, products._Description, products._Fld37,
                   products._Fld38
          having sum(balance._Fld249) > 0

          union all

          select data_res.collection,
                 data_res.uuid,
                 data_res.code,
                 data_res.description,
                 data_res.length,
                 data_res.width,
                 data_res.count
          from data_res)
select data_collection.collection,
       data_collection.uuid,
       data_collection.code,
       data_collection.description,
       data_collection.length,
       data_collection.width,
       sum(data_collection.count) as count
from data_collection
group by data_collection.collection, data_collection.uuid, data_collection.code, data_collection.description,
         data_collection.length, data_collection.width
having sum(data_collection.count) > 0