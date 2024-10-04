with data_res as (select products.collection          as collection,
                         encode(products.uuid, 'hex') as uuid,
                         products.code                as code,
                         products.description         as description,
                         products.length              as length,
                         products.width               as width,
                         -sum(reserved.reserved)      as count
                  from _reserved as reserved
                           inner join _products as products
                                      on reserved.productUUID = products.uuid
                  where true
                  group by products.collection, products.uuid, products.code, products.description, products.length,
                           products.width, reserved.detailsUUID
                  having sum(reserved.reserved) > 0),

     data_collection as
         (select products.collection          as collection,
                 encode(products.uuid, 'hex') as uuid,
                 products.code                as code,
                 products.description         as description,
                 products.length              as length,
                 products.width               as width,
                 sum(balance.remaining)       as count
          from _balance as balance
                   inner join _products as products
                              on balance.productUUID = products.uuid
          where true
          group by products.collection, products.uuid, products.code, products.description, products.length,
                   products.width
          having sum(balance.remaining) > 0

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