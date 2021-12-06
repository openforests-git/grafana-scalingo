import React from 'react';
import { cx } from '@emotion/css';

export interface CardProps {
  logoUrl?: string;
  logoAlt?: string;
  title: string;
  description?: string;
  labels?: React.ReactNode;
  actions?: React.ReactNode;
  onClick?: () => void;
  ariaLabel?: string;
  className?: string;
}

export const Card: React.FC<CardProps> = ({
  logoUrl,
  logoAlt,
  title,
  description,
  labels,
  actions,
  onClick,
  ariaLabel,
  className,
}) => {
  const mainClassName = cx('add-data-source-item', className);

  return (
    <div className={mainClassName} onClick={onClick} aria-label={ariaLabel}>
      {logoUrl && <img className="add-data-source-item-logo" src={logoUrl} alt={logoAlt ?? ''} />}
      <div className="add-data-source-item-text-wrapper">
        <span className="add-data-source-item-text">{title}</span>
        {description && <span className="add-data-source-item-desc">{description}</span>}
        {labels && <div className="add-data-source-item-badge">{labels}</div>}
      </div>
      {actions && <div className="add-data-source-item-actions">{actions}</div>}
    </div>
  );
};
