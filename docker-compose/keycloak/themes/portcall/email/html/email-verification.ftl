<#import "template.ftl" as layout>
<@layout.emailLayout>
<h1>Verify Your Email</h1>
<p>Someone has created a <strong>${realmName}</strong> account with this email address.</p>
<p>If this was you, click the button below to verify your email address:</p>
<p style="text-align: center;">
    <a href="${link}" class="button">Verify Email</a>
</p>
<p>Or copy and paste this link into your browser:</p>
<p style="word-break: break-all; background-color: #f4f4f5; padding: 12px; border-radius: 6px; font-size: 14px;">
    <a href="${link}">${link}</a>
</p>
<p>This link will expire within <strong>${linkExpirationFormatter(linkExpiration)}</strong>.</p>
<p style="color: #71717a; font-size: 14px; margin-top: 24px;">If you didn't create this account, you can safely ignore this email.</p>
</@layout.emailLayout>
