with data_dtls as
         (select products.collection                as collection,
                 encode(products.uuid, 'hex')       as uuid,
                 products.code                      as code,
                 products.description               as description,
                 products.length                    as length,
                 products.width                     as width,
                 encode(details.uuidDetails, 'hex') as uuidDetails,
                 details.codeDetails                as codeDetails,
                 details.color                      as color,
                 details.picture                    as picture,
                 details.form                       as form,
                 details.barcode                    as barcode,
                 coalesce(brands.brandName, '')     as brand,
                 sum(balance.countBalance)          as countBalance,
                 0                                  as reservedBalance,
                 sum(balance.countBalance)          as count
          from _balance as balance
                   inner join _products as products
                              on balance.productUUID = products.uuid
                   inner join _details as details
                              on balance.detailsUUID = details.uuidDetails
                   left join _brands as brands
                             on details.brandUUID = brands.brandUUID
          where true
          group by products.collection, encode(products.uuid, 'hex'), products.code, products.description,
                   products.length,
                   products.width, encode(details.uuidDetails, 'hex'), details.codeDetails, details.color,
                   details.picture,
                   details.form, details.barcode, coalesce(brands.brandName, '')
          having sum(balance.countBalance) > 0

          union all

          select products.collection,
                 encode(products.uuid, 'hex'),
                 products.code,
                 products.description,
                 products.length,
                 products.width,
                 encode(details.uuidDetails, 'hex'),
                 details.codeDetails,
                 details.color,
                 details.picture,
                 details.form,
                 details.barcode,
                 coalesce(brands.brandName, ''),
                 0,
                 sum(reserved.reservedBalance),
                 -sum(reserved.reservedBalance)
          from _reserved as reserved
                   inner join _products as products
                              on reserved.productUUID = products.uuid
                   inner join _details as details
                              on reserved.detailsUUID = details.uuidDetails
                   left join _brands as brands
                             on details.brandUUID = brands.brandUUID
          where true
          group by products.collection, encode(products.uuid, 'hex'), products.code, products.description,
                   products.length,
                   products.width, encode(details.uuidDetails, 'hex'), details.codeDetails, details.color,
                   details.picture,
                   details.form, details.barcode, coalesce(brands.brandName, '')
          having sum(reserved.reservedBalance) > 0)
select data_dtls.collection,
       data_dtls.uuid,
       data_dtls.code,
       data_dtls.description,
       data_dtls.length,
       data_dtls.width,
       data_dtls.uuidDetails,
       data_dtls.codeDetails,
       data_dtls.color,
       data_dtls.picture,
       data_dtls.form,
       data_dtls.barcode,
       data_dtls.brand,
       sum(data_dtls.countBalance)    as countBalance,
       sum(data_dtls.reservedBalance) as reservedBalance,
       sum(data_dtls.count)           as count
from data_dtls
group by data_dtls.collection, data_dtls.uuid, data_dtls.code, data_dtls.description, data_dtls.length,
         data_dtls.width, data_dtls.uuidDetails, data_dtls.codeDetails, data_dtls.color,
         data_dtls.picture, data_dtls.form, data_dtls.barcode, data_dtls.brand
having sum(data_dtls.count) > 0
