import React from 'react';
import {Dropdown} from './Dropdown';

export default ({
    items,
    onChange,
    trigger,
    triggerClassName
}) => (
    <Dropdown
        trigger={trigger}
        disabled={false}
        triggerClassName={ triggerClassName }
        isClosableOnContentClick={false}
    >
        <ul className={'DropdownList-list'}>
            { items.map((item, key) => (
                <li
                    key={key}
                    className={'DropdownList-listItem'}
                    onClick={()=>onChange(item.value)}
                >
                    {item.text}
                </li>
            ))}
        </ul>
    </Dropdown>
);
