import React from 'react';
import './Input.css';

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
    label?: string;
}

export const Input: React.FC<InputProps> = ({ label, ...props }) => {
    return (
        <div className="ui-input-wrapper">
            {label && <label className="ui-input-label">{label}</label>}
            <input className="ui-input" {...props} />
        </div>
    );
};
