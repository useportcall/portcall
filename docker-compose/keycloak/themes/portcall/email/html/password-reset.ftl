<#import "template.ftl" as layout>
<@layout.emailLayout>
<h1>Reset Your Password</h1>
<p>Someone just requested to reset your password for your <strong>${realmName}</strong> account.</p>
<p>If this was you, click the button below to reset your password:</p>
<p style="text-align: center;">
    <a href="${link}" class="button">Reset Password</a>
</p>
<p>Or copy and paste this link into your browser:</p>
<p style="word-break: break-all; background-color: #f4f4f5; padding: 12px; border-radius: 6px; font-size: 14px;">
    <a href="${link}">${link}</a>
</p>
<p>This link will expire within <strong>${linkExpirationFormatter(linkExpiration)}</strong>.</p>
<p style="color: #71717a; font-size: 14px; margin-top: 24px;">If you didn't request a password reset, you can safely ignore this email. Your password will remain unchanged.</p>
</@layout.emailLayout>
