select img.imageID                as imageID,
       img.collection             as collection,
       img.picture                as picture,
       img.form                   as form,
       img.color                 as color,
--        brands.brandCode as brandCode,
       coalesce(brands.brand, '') as brand
from _img as img
         left join _brands as brands
                   on img.brandUUID = brands.uuid
where true