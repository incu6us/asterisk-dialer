import React from 'react';

export const Button = ({
   className,
   onClick,
   inscription,
   isDisabled
}) => {
    return (
        <button
            className={className}
            onClick={onClick}
            disabled={isDisabled}>
            {inscription}
        </button>
    )
};
