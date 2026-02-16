<#import "template.ftl" as layout>
<@layout.emailLayout>
<h1>Verify your email</h1>
<p>A <strong style="color: #0a0a0a;">${realmName}</strong> account was created with this email address.</p>
<p>Click the button below to verify your email:</p>
<p style="text-align: center; padding: 8px 0;">
    <a href="${link}" class="button">Verify email</a>
</p>
<p>Or copy and paste this link into your browser:</p>
<p class="link-box">
    <a href="${link}">${link}</a>
</p>
<p>This link expires in <strong style="color: #0a0a0a;">${linkExpirationFormatter(linkExpiration)}</strong>.</p>
<p class="hint">If you didn't create this account, you can safely ignore this email.</p>
</@layout.emailLayout>
