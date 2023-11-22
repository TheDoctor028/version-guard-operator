create table DeploymentChangeData
(
    kind varchar not null,
    name varchar not null,
    namespace varchar not null,
    selector varchar not null,
    container_name varchar not null,
    image varchar not null,
    ts timestamp not null
);

create table  ApplicationChangeData
(
   kind varchar not null,
   name varchar not null,
   namespace varchar not null,
   selector varchar not null,
   container_name varchar not null,
   image varchar not null,
   ts timestamp not null
);

alert table DeploymentChangeData add constraint PK_DeploymentChangeData primary key (name, namespace, selector);
alert table ApplicationChangeData add constraint PK_ApplicationChangeData primary key (name, namespace, selector);


create view MismatchingVersions
            (
             kind,
             name,
             namespace,
             container_name,
             current_image,
             desired_image,
             ts
                ) as
select
    DCD.kind as kind,
    ACD.name as name,
    ACD.namespace as namespace,
    ACD.container_name as container_name,
    DCD.image as current_image,
    ACD.image as desired_image,
    DCD.ts as ts
from
    DeploymentChangeData as DCD, ApplicationChangeData as ACD
where
        DCD.namespace = ACD.namespace and
        DCD.selector = ACD.selector and
        DCD.image != ACD.image;
