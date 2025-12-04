# Hospital EMR System - Security & Compliance Guide

## Overview

This document outlines the security measures and compliance standards implemented in the Hospital EMR System to meet HIPAA, UU PDP Indonesia, and other healthcare data protection requirements.

---

## Table of Contents

1. [Security Architecture](#security-architecture)
2. [HIPAA Compliance](#hipaa-compliance)
3. [UU PDP Compliance](#uu-pdp-compliance)
4. [Data Encryption](#data-encryption)
5. [Access Control](#access-control)
6. [Audit Logging](#audit-logging)
7. [Network Security](#network-security)
8. [Incident Response](#incident-response)

---

## Security Architecture

### Zero Trust Model

The system implements a **Zero Trust** security model where:
- No implicit trust is granted to any user or system
- All access must be authenticated and authorized
- Continuous verification of security posture
- Least privilege access principle

### Defense in Depth

Multiple layers of security controls:
1. **Network Layer**: Firewalls, VPCs, security groups
2. **Application Layer**: Authentication, authorization, input validation
3. **Data Layer**: Encryption, access controls, audit logging
4. **Infrastructure Layer**: Container security, K8s policies

---

## HIPAA Compliance

### HIPAA Security Rule

#### Administrative Safeguards

✅ **Security Management Process**
- Risk analysis conducted
- Risk management strategy implemented
- Sanction policy for violations
- Information system activity review (audit logs)

✅ **Assigned Security Responsibility**
- Designated Security Officer role
- Clear security responsibilities documented

✅ **Workforce Security**
- Authorization and supervision procedures
- Workforce clearance procedures
- Termination procedures

✅ **Information Access Management**
- Role-Based Access Control (RBAC)
- Minimum necessary access principle
- Access authorization procedures

✅ **Security Awareness and Training**
- Security reminders
- Protection from malicious software
- Login monitoring
- Password management

✅ **Security Incident Procedures**
- Response and reporting procedures documented
- Incident response team established

✅ **Contingency Plan**
- Data backup plan (automated daily backups)
- Disaster recovery plan (RPO < 1 hour, RTO < 4 hours)
- Emergency mode operation plan
- Testing and revision procedures

✅ **Business Associate Agreements**
- BAAs with all third-party vendors (Neon, cloud providers)

#### Physical Safeguards

✅ **Facility Access Controls**
- Cloud infrastructure with physical security
- Visitor control and authorization
- Access control and validation procedures

✅ **Workstation Security**
- Workstation use policies
- Automatic logout after inactivity

✅ **Device and Media Controls**
- Disposal procedures (secure data deletion)
- Media re-use procedures
- Backup and storage procedures

#### Technical Safeguards

✅ **Access Control**
- Unique user identification (UUIDs)
- Emergency access procedure
- Automatic logoff (30-minute timeout)
- Encryption and decryption (AES-256)

✅ **Audit Controls**
- Comprehensive audit trails for all PHI access
- Immutable logs stored for 25 years
- Log includes: WHO, WHAT, WHEN, WHERE

✅ **Integrity**
- Hash verification for data integrity
- Digital signatures for critical documents

✅ **Person or Entity Authentication**
- Multi-factor authentication (FIDO2/WebAuthn)
- Strong password requirements
- Session management

✅ **Transmission Security**
- TLS 1.3 for all communications
- mTLS for service-to-service communication

### HIPAA Privacy Rule

✅ **Patient Rights**
- Right to access PHI (Patient Portal)
- Right to amend PHI
- Right to accounting of disclosures
- Right to request restrictions
- Right to confidential communications

✅ **Notice of Privacy Practices**
- Patients receive notice upon registration
- Available on patient portal

✅ **Minimum Necessary Standard**
- Row-Level Security (RLS) enforces minimum necessary access
- Role-based permissions limit data access

---

## UU PDP Compliance

### Indonesia Personal Data Protection Law (UU No. 27/2022)

✅ **Consent Management**
- Explicit consent collected for data processing
- Consent withdrawal mechanism
- Purpose-specific consent

✅ **Data Subject Rights**
- Right to access personal data
- Right to rectification
- Right to erasure ("right to be forgotten")
- Right to data portability
- Right to object to processing

✅ **Data Processing Principles**
- Lawfulness, fairness, transparency
- Purpose limitation
- Data minimization
- Accuracy
- Storage limitation
- Integrity and confidentiality

✅ **Data Protection Officer**
- DPO designated and contactable

✅ **Data Breach Notification**
- Breach detection mechanisms
- 72-hour notification requirement
- Incident response procedures

### PMK No. 24/2022 (Electronic Medical Records)

✅ **EMR Retention**
- Medical records stored for minimum 25 years
- Automated retention policy
- Secure archival procedures

✅ **Data Integrity**
- Digital signatures for medical documents
- Audit trails prevent tampering
- Version control for document changes

✅ **Confidentiality**
- Encryption at rest and in transit
- Access controls and authentication
- Confidentiality agreements with staff

---

## Data Encryption

### Encryption at Rest

**AES-256 Encryption**
- All PHI encrypted in database
- Full disk encryption on servers
- Encrypted backups

**Key Management**
- Keys stored in secure Key Management Service (KMS)
- Key rotation every 90 days
- Separation of keys from encrypted data

**Implementation**:
```go
// Sensitive fields encrypted before storage
encryptedSSN, err := encryption.Encrypt(patient.SSN, encryptionKey)
```

### Encryption in Transit

**TLS 1.3**
- All client-server communication
- Minimum TLS version enforced
- Strong cipher suites only

**mTLS (Mutual TLS)**
- Service-to-service communication
- Certificate-based authentication
- Enforced by service mesh (Istio/Linkerd)

**Database Connections**
- TLS required for all database connections
- Certificate verification: `sslmode=verify-full`

---

## Access Control

### Authentication

**Multi-Factor Authentication (MFA)**
- FIDO2/WebAuthn support
- TOTP (Time-based One-Time Password)
- Required for administrative access
- Optional but recommended for all users

**Password Requirements**
- Minimum 12 characters
- Complexity requirements (uppercase, lowercase, numbers, symbols)
- Password history (last 10 passwords)
- Maximum age: 90 days
- Account lockout after 5 failed attempts

**Session Management**
- JWT-based authentication
- Session timeout: 30 minutes of inactivity
- Secure session storage
- Session revocation on logout

### Authorization

**Role-Based Access Control (RBAC)**

Predefined roles:
- **Administrator**: Full system access
- **Doctor**: Clinical data access, order entry, prescriptions
- **Nurse**: Patient care data, vital signs, medication administration
- **Receptionist**: Patient registration, appointment scheduling
- **Lab Technician**: Lab orders and results
- **Radiologist**: Radiology orders and results
- **Pharmacist**: Prescription access and dispensing

**Permissions**:
```
view_patients
create_patients
update_patients
delete_patients
view_encounters
create_encounters
update_encounters
view_orders
create_orders
view_results
update_results
manage_users
manage_roles
view_audit_log
```

**Row-Level Security (RLS)**
- Database policies restrict data access
- Doctors only see their patients
- Patients only see their own data
- Cross-referencing based on care relationships

**Example RLS Policy**:
```sql
CREATE POLICY doctor_patients_policy ON patients
    FOR SELECT
    USING (
        EXISTS (
            SELECT 1 FROM encounters
            WHERE encounters.patient_id = patients.id
            AND encounters.provider_id = current_user_id()
        )
    );
```

### API Security

**Rate Limiting**
- 100 requests per minute per user
- Prevents brute force attacks
- DDoS protection

**Input Validation**
- All inputs validated and sanitized
- SQL injection prevention (parameterized queries)
- XSS prevention (output encoding)
- CSRF protection

---

## Audit Logging

### What is Logged

Every action involving PHI is logged:
- **CREATE**: New record creation
- **READ**: Data access/viewing
- **UPDATE**: Data modification
- **DELETE**: Data deletion
- **PRINT**: Document printing
- **EXPORT**: Data export
- **LOGIN/LOGOUT**: User sessions

### Log Content

Each audit log entry contains:
- **Timestamp**: UTC timestamp of action
- **User ID**: Who performed the action
- **Username**: User's email/identifier
- **Action**: Type of action (CREATE, READ, UPDATE, DELETE)
- **Resource**: Type of resource (patient, encounter, order)
- **Resource ID**: Specific record ID
- **IP Address**: Origin of request
- **User Agent**: Browser/client information
- **Changes**: Before/after values for updates
- **Status**: Success or failure

### Log Retention

- **Duration**: 25 years (compliance requirement)
- **Storage**: Immutable storage (WORM - Write Once, Read Many)
- **Access**: Restricted to authorized personnel only
- **Review**: Regular audit log reviews for anomalies

### Example Audit Log

```json
{
  "id": "uuid",
  "timestamp": "2025-01-06T10:30:00Z",
  "user_id": "uuid",
  "username": "doctor@hospital-emr.com",
  "action": "READ",
  "resource": "patient",
  "resource_id": "patient-uuid",
  "description": "Viewed patient medical record",
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "status_code": 200,
  "severity": "info"
}
```

---

## Network Security

### Firewall Configuration

**Inbound Rules**:
- Port 443 (HTTPS): Public access to API
- Port 22 (SSH): Restricted to admin IPs only
- Port 5432 (PostgreSQL): Internal network only
- Port 6379 (Redis): Internal network only
- Port 4222 (NATS): Internal network only

**Outbound Rules**:
- Allow HTTPS (443) for external API calls
- Allow DNS (53)
- Deny all others by default

### VPC and Subnets

**Network Segmentation**:
- Public subnet: Load balancers only
- Private subnet: Application servers
- Data subnet: Databases (no internet access)

**Security Groups**:
- Least privilege principle
- No unnecessary ports open
- Regular security group audits

### DDoS Protection

- CloudFlare or AWS Shield
- Rate limiting at API Gateway
- Geographic restrictions if applicable

---

## Incident Response

### Incident Response Plan

**Phases**:
1. **Preparation**: Team, tools, procedures ready
2. **Detection**: Monitoring and alerting
3. **Containment**: Isolate affected systems
4. **Eradication**: Remove threat
5. **Recovery**: Restore services
6. **Lessons Learned**: Post-incident review

### Breach Notification

**Timeline**:
- Internal notification: Immediate
- Security team assessment: Within 4 hours
- Management notification: Within 24 hours
- Regulatory notification: Within 72 hours (if required)
- Affected individuals: Without unreasonable delay

**Notification Content**:
- Description of breach
- Types of data involved
- Steps taken to mitigate
- Contact information
- Steps individuals should take

### Contact Information

- **Security Team**: security@hospital-emr.com
- **Emergency Hotline**: +1-XXX-XXX-XXXX
- **Incident Report**: https://security.hospital-emr.com/report

---

## Security Best Practices

### For Developers

1. Never commit secrets to version control
2. Use parameterized queries (prevent SQL injection)
3. Validate and sanitize all inputs
4. Follow principle of least privilege
5. Conduct code reviews focusing on security
6. Keep dependencies up to date
7. Run security scans (SAST/DAST) before deployment

### For Administrators

1. Enable MFA for all accounts
2. Use strong, unique passwords
3. Regularly review access logs
4. Keep systems patched and updated
5. Conduct regular security audits
6. Test backup and recovery procedures
7. Monitor for suspicious activity

### For End Users

1. Use strong passwords
2. Enable MFA
3. Don't share credentials
4. Log out when finished
5. Report suspicious activity
6. Attend security training
7. Follow data handling policies

---

## Compliance Checklist

### Pre-Production

- [ ] Security risk assessment completed
- [ ] Penetration testing performed
- [ ] HIPAA compliance audit passed
- [ ] UU PDP compliance verified
- [ ] BAAs signed with all vendors
- [ ] Disaster recovery plan tested
- [ ] Incident response plan documented
- [ ] Security training completed
- [ ] Encryption verified (at-rest and in-transit)
- [ ] Access controls tested
- [ ] Audit logging verified
- [ ] Backup and restore tested

### Ongoing

- [ ] Monthly security log reviews
- [ ] Quarterly access audits
- [ ] Semi-annual penetration testing
- [ ] Annual HIPAA audit
- [ ] Continuous vulnerability scanning
- [ ] Regular security training
- [ ] Patch management process
- [ ] Backup verification

---

## Certifications

Target certifications:
- **HIPAA Compliance**: In progress
- **ISO 27001**: Planned
- **SOC 2 Type II**: Planned
- **HITRUST**: Future consideration

---

## References

- HIPAA Security Rule: https://www.hhs.gov/hipaa/for-professionals/security/
- HIPAA Privacy Rule: https://www.hhs.gov/hipaa/for-professionals/privacy/
- UU No. 27 Tahun 2022: Indonesian Personal Data Protection Law
- PMK No. 24 Tahun 2022: Electronic Medical Records Regulation
- NIST Cybersecurity Framework: https://www.nist.gov/cyberframework
- OWASP Top 10: https://owasp.org/www-project-top-ten/

---

## Contact

For security concerns or to report vulnerabilities:
- Email: security@hospital-emr.com
- Bug Bounty: https://bugcrowd.com/hospital-emr
- PGP Key: Available on request

---

*Last Updated: 2025-01-06*
*Version: 1.0*
