with data_dtls as
         (select products._Fld36                as collection,
                 encode(products._IDRRef, 'hex')       as uuid,
                 products._Code                      as code,
                 products._Description               as description,
                 products._Fld37                    as length,
                 products._Fld38                     as width,
                 encode(details._IDRRef, 'hex') as uuidDetails,
                 details._Code                as codeDetails,
                 details._Fld60                      as color,
                 details._Fld61                    as picture,
                 details._Fld62                       as form,
                 details._Fld63                    as barcode,
                 coalesce(brands._Description, '')     as brand,
                 sum(balance._Fld249)          as countBalance,
                 0                                  as reservedBalance,
                 sum(balance._Fld249)          as count
          from _AccumRgT251 as balance
                   inner join _Reference12 as products
                              on balance._Fld246RRef = products._IDRRef
                   inner join _Reference16 as details
                              on balance._Fld247RRef = details._IDRRef
              --          left join _AccumRgT244 as reserved
--                    on balance._Fld246RRef = reserved._Fld239RRef and balance._Fld247RRef = reserved._Fld240RRef
                   left join _Reference13 as brands
                             on details._Fld59RRef = brands._IDRRef
          where true
          group by products._Fld36, encode(products._IDRRef, 'hex'), products._Code, products._Description,
                   products._Fld37,
                   products._Fld38, encode(details._IDRRef, 'hex'), details._Code, details._Fld60,
                   details._Fld61,
                   details._Fld62, details._Fld63, coalesce(brands._Description, '')
          having sum(balance._Fld249) > 0

          union all

          select products._Fld36,
                 encode(products._IDRRef, 'hex'),
                 products._Code,
                 products._Description,
                 products._Fld37,
                 products._Fld38,
                 encode(details._IDRRef, 'hex'),
                 details._Code,
                 details._Fld60,
                 details._Fld61,
                 details._Fld62,
                 details._Fld63,
                 coalesce(brands._Description, ''),
                 0,
                 sum(reserved._Fld243),
                 -sum(reserved._Fld243)
          from _AccumRgT244 as reserved
                   inner join _Reference12 as products
                              on reserved._Fld239RRef = products._IDRRef
                   inner join _Reference16 as details
                              on reserved._Fld240RRef = details._IDRRef
                   left join _Reference13 as brands
                             on details._Fld59RRef = brands._IDRRef
          where true
          group by products._Fld36, encode(products._IDRRef, 'hex'), products._Code, products._Description,
                   products._Fld37,
                   products._Fld38, encode(details._IDRRef, 'hex'), details._Code, details._Fld60,
                   details._Fld61,
                   details._Fld62, details._Fld63, coalesce(brands._Description, '')
          having sum(reserved._Fld243) > 0)
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
